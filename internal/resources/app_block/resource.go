// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/path"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/metadata"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/tags"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
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

	// determine packaging type (aws defaults to CUSTOM if unset)
	packagingType := "CUSTOM"
	if !config.PackagingType.IsNull() && !config.PackagingType.IsUnknown() {
		packagingType = config.PackagingType.ValueString()
	}

	hasSetupScript := !config.SetupScriptDetails.IsNull() && !config.SetupScriptDetails.IsUnknown()
	hasPostSetupScript := !config.PostSetupScriptDetails.IsNull() && !config.PostSetupScriptDetails.IsUnknown()

	switch packagingType {

	case "CUSTOM":
		// setup script is mandatory
		if !hasSetupScript {
			resp.Diagnostics.AddAttributeError(
				path.Root("setup_script_details"),
				"Missing Required Attribute",
				"`setup_script_details` must be specified when `packaging_type` is `CUSTOM`.",
			)
		}

		// post-setup script not allowed
		if hasPostSetupScript {
			resp.Diagnostics.AddAttributeError(
				path.Root("post_setup_script_details"),
				"Invalid Configuration",
				"`post_setup_script_details` can only be specified when `packaging_type` is `APPSTREAM2`.",
			)
		}

	case "APPSTREAM2":
		// setup script not allowed
		if hasSetupScript {
			resp.Diagnostics.AddAttributeError(
				path.Root("setup_script_details"),
				"Invalid Configuration",
				"`setup_script_details` cannot be specified when `packaging_type` is `APPSTREAM2`.",
			)
		}
	}

	// validate source_s3_location.s3_key requirements
	if !config.SourceS3Location.IsNull() && !config.SourceS3Location.IsUnknown() {
		var source sourceS3LocationModel

		resp.Diagnostics.Append(
			config.SourceS3Location.As(ctx, &source, basetypes.ObjectAsOptions{})...,
		)
		if resp.Diagnostics.HasError() {
			return
		}

		hasS3Key := !source.S3Key.IsNull() && !source.S3Key.IsUnknown()

		if packagingType == "CUSTOM" && !hasS3Key {
			resp.Diagnostics.AddAttributeError(
				path.Root("source_s3_location").AtName("s3_key"),
				"Missing Required Attribute",
				"`s3_key` must be specified in `source_s3_location` when `packaging_type` is `CUSTOM`.",
			)
		}
	}
}

func (r *resource) Metadata(_ context.Context, req tfresource.MetadataRequest, resp *tfresource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_block"
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
	if err := util.ValidateARNValue(req.ID, "appstream", "app-block/"); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected AppStream app block ARN: %v", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
