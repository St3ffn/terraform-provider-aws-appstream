// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/path"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
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

	if !config.StorageConnectors.IsNull() && !config.StorageConnectors.IsUnknown() {
		var connectors []storageConnectorModel
		resp.Diagnostics.Append(
			config.StorageConnectors.ElementsAs(ctx, &connectors, false)...,
		)
		if resp.Diagnostics.HasError() {
			return
		}

		for i, c := range connectors {
			if c.ConnectorType.IsUnknown() || c.DomainsRequireAdminConsent.IsUnknown() {
				continue
			}

			if !c.DomainsRequireAdminConsent.IsNull() &&
				c.ConnectorType.ValueString() != "ONE_DRIVE" {

				resp.Diagnostics.AddAttributeError(
					path.Root("storage_connectors").AtListIndex(i).AtName("domains_require_admin_consent"),
					"Invalid Storage Connector Configuration",
					"`domains_require_admin_consent` can only be specified when `connector_type` is `ONE_DRIVE`.",
				)
			}
		}
	}

	if !config.UserSettings.IsNull() && !config.UserSettings.IsUnknown() {
		var settings []userSettingModel
		resp.Diagnostics.Append(
			config.UserSettings.ElementsAs(ctx, &settings, false)...,
		)
		if resp.Diagnostics.HasError() {
			return
		}

		for i, s := range settings {
			if s.Action.IsUnknown() || s.Permission.IsUnknown() || s.MaximumLength.IsUnknown() {
				continue
			}

			if !s.MaximumLength.IsNull() {
				action := s.Action.ValueString()
				permission := s.Permission.ValueString()

				if action != "CLIPBOARD_COPY_FROM_LOCAL_DEVICE" &&
					action != "CLIPBOARD_COPY_TO_LOCAL_DEVICE" {

					resp.Diagnostics.AddAttributeError(
						path.Root("user_settings").AtListIndex(i).AtName("maximum_length"),
						"Invalid User Setting",
						"`maximum_length` can only be specified for clipboard copy actions.",
					)
				}

				if permission == "DISABLED" {
					resp.Diagnostics.AddAttributeError(
						path.Root("user_settings").AtListIndex(i).AtName("maximum_length"),
						"Invalid User Setting",
						"`maximum_length` cannot be specified when `permission` is `DISABLED`.",
					)
				}
			}
		}
	}
}

func (r *resource) Metadata(_ context.Context, req tfresource.MetadataRequest, resp *tfresource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack"
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
			"Expected import identifier format: <stack_name>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
