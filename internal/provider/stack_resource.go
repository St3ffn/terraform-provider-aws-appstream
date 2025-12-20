// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &stackResource{}
	_ resource.ResourceWithConfigure      = &stackResource{}
	_ resource.ResourceWithValidateConfig = &stackResource{}
	_ resource.ResourceWithImportState    = &stackResource{}
)

func NewStackResource() resource.Resource {
	return &stackResource{}
}

type stackResource struct {
	appstreamClient *awsappstream.Client
	taggingClient   *awstaggingapi.Client
}

func (r *stackResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config stackModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.StorageConnectors.IsNull() && !config.StorageConnectors.IsUnknown() {
		var connectors []stackStorageConnectorModel
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
		var settings []stackUserSettingModel
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

func (r *stackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack"
}

func (r *stackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*metadata)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *metadata, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if meta.appstream == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *metadata.appstream, got: nil. Please report this issue to the provider developers.",
		)
		return
	}

	if meta.tagging == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *metadata.tagging, got: nil. Please report this issue to the provider developers.",
		)
		return
	}

	r.appstreamClient = meta.appstream
	r.taggingClient = meta.tagging
}

func (r *stackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
