// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/metadata"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/tags"
)

var (
	_ tfresource.Resource                   = &resource{}
	_ tfresource.ResourceWithConfigure      = &resource{}
	_ tfresource.ResourceWithValidateConfig = &resource{}
	_ tfresource.ResourceWithImportState    = &resource{}
)

func NewResource() tfresource.Resource {
	return &resource{}
}

type resource struct {
	appstreamClient *awsappstream.Client
	tags            *tags.TagManager
}

func (r *resource) ValidateConfig(ctx context.Context, req tfresource.ValidateConfigRequest, resp *tfresource.ValidateConfigResponse) {
	var config model

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasImageName := !config.ImageName.IsNull() && !config.ImageName.IsUnknown()
	hasImageArn := !config.ImageARN.IsNull() && !config.ImageARN.IsUnknown()

	// can either have image name or image arn
	switch {
	case hasImageName && hasImageArn:
		resp.Diagnostics.AddAttributeError(
			path.Root("image_name"),
			"Invalid Image Configuration",
			"Only one of `image_name` or `image_arn` may be specified.",
		)

	case !hasImageName && !hasImageArn:
		resp.Diagnostics.AddAttributeError(
			path.Root("image_name"),
			"Missing Required Attribute",
			"Either `image_name` or `image_arn` must be specified.",
		)
	}

	fleetType := config.FleetType.ValueString()

	switch fleetType {

	case string(awstypes.FleetTypeElastic):
		// elastic fleets scale by sessions, not instances
		if config.MaxConcurrentSessions.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("max_concurrent_sessions"),
				"Missing Required Attribute",
				"`max_concurrent_sessions` must be specified for elastic fleets.",
			)
		}

		// per-instance session limits are not applicable to elastic fleets
		if !config.MaxSessionsPerInstance.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("max_sessions_per_instance"),
				"Invalid Configuration",
				"`max_sessions_per_instance` cannot be specified for elastic fleets.",
			)
		}

		// elastic fleets do not use compute capacity blocks
		if !config.ComputeCapacity.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("compute_capacity"),
				"Invalid Configuration",
				"`compute_capacity` cannot be specified for elastic fleets.",
			)
		}

		// elastic fleets always require VPC configuration
		if config.VPCConfig.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vpc_config"),
				"Missing Required Attribute",
				"`vpc_config` must be specified for elastic fleets.",
			)
		}

		var vpc vpcConfigModel
		resp.Diagnostics.Append(
			config.VPCConfig.As(ctx, &vpc, basetypes.ObjectAsOptions{})...,
		)
		if resp.Diagnostics.HasError() {
			return
		}

		// subnets are mandatory for elastic fleets
		if vpc.SubnetIDs.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vpc_config").AtName("subnet_ids"),
				"Missing Required Attribute",
				"`subnet_ids` must be specified for elastic fleets.",
			)
		} else if !vpc.SubnetIDs.IsUnknown() {
			// aws requires at least two subnets in different AZs
			subnets := vpc.SubnetIDs.Elements()
			if len(subnets) < 2 {
				resp.Diagnostics.AddAttributeError(
					path.Root("vpc_config").AtName("subnet_ids"),
					"Invalid VPC Configuration",
					"Elastic fleets require at least two subnets in different availability zones.",
				)
			}
		}

		// domain join is not supported for elastic fleets
		if !config.DomainJoinInfo.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("domain_join_info"),
				"Invalid Configuration",
				"`domain_join_info` cannot be specified for elastic fleets.",
			)
		}

	case string(awstypes.FleetTypeOnDemand), string(awstypes.FleetTypeAlwaysOn):
		// session-based scaling is exclusive to elastic fleets
		if !config.MaxConcurrentSessions.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("max_concurrent_sessions"),
				"Invalid Configuration",
				"`max_concurrent_sessions` can only be specified for elastic fleets.",
			)
		}

		// session scripts are only supported by elastic fleets
		if !config.SessionScriptS3Location.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("session_script_s3_location"),
				"Invalid Configuration",
				"`session_script_s3_location` can only be specified for elastic fleets.",
			)
		}

		// non-elastic fleets must define compute capacity
		if config.ComputeCapacity.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("compute_capacity"),
				"Missing Required Attribute",
				"`compute_capacity` must be specified for non-elastic fleets.",
			)
			return
		}

		var capacity computeCapacityModel
		resp.Diagnostics.Append(
			config.ComputeCapacity.As(ctx, &capacity, basetypes.ObjectAsOptions{})...,
		)
		if resp.Diagnostics.HasError() {
			return
		}

		hasDesiredInstances := !capacity.DesiredInstances.IsNull() && !capacity.DesiredInstances.IsUnknown()
		hasDesiredSessions := !capacity.DesiredSessions.IsNull() && !capacity.DesiredSessions.IsUnknown()

		switch {
		// aws requires exactly one of instances or sessions
		case hasDesiredInstances && hasDesiredSessions:
			resp.Diagnostics.AddAttributeError(
				path.Root("compute_capacity"),
				"Invalid Compute Capacity Configuration",
				"Only one of `desired_instances` or `desired_sessions` may be specified.",
			)

		case !hasDesiredInstances && !hasDesiredSessions:
			resp.Diagnostics.AddAttributeError(
				path.Root("compute_capacity"),
				"Missing Required Attribute",
				"Exactly one of `desired_instances` or `desired_sessions` must be specified for non-elastic fleets.",
			)

		// desired_sessions implies a multi-session fleet
		case hasDesiredSessions:
			if config.MaxSessionsPerInstance.IsNull() || config.MaxSessionsPerInstance.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("compute_capacity").AtName("desired_sessions"),
					"Invalid Compute Capacity Configuration",
					"`desired_sessions` can only be used for multi-session fleets. "+
						"Set `max_sessions_per_instance` to a value greater than 1.",
				)
			} else if config.MaxSessionsPerInstance.ValueInt32() <= 1 {
				resp.Diagnostics.AddAttributeError(
					path.Root("compute_capacity").AtName("desired_sessions"),
					"Invalid Compute Capacity Configuration",
					"`desired_sessions` requires `max_sessions_per_instance` to be greater than 1.",
				)
			}
		}
	}
}

func (r *resource) Metadata(_ context.Context, req tfresource.MetadataRequest, resp *tfresource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fleet"
}

func (r *resource) Configure(_ context.Context, req tfresource.ConfigureRequest, resp *tfresource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*metadata.Metadata)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Metadata, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if meta.Appstream == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *Metadata.Appstream, got: nil. Please report this issue to the provider developers.",
		)
		return
	}

	if meta.Tagging == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *Metadata.Tagging, got: nil. Please report this issue to the provider developers.",
		)
		return
	}

	r.appstreamClient = meta.Appstream
	r.tags = tags.NewTagManager(meta.Tagging, meta.DefaultTags)
}

func (r *resource) ImportState(ctx context.Context, req tfresource.ImportStateRequest, resp *tfresource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier format: <fleet_name>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
