// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var AppStreamMaxResults int32 = 50

var (
	_ resource.Resource                = &associateApplicationEntitlementResource{}
	_ resource.ResourceWithConfigure   = &associateApplicationEntitlementResource{}
	_ resource.ResourceWithImportState = &associateApplicationEntitlementResource{}
)

func NewAssociateApplicationEntitlementResource() resource.Resource {
	return &associateApplicationEntitlementResource{}
}

type associateApplicationEntitlementResource struct {
	appstreamClient *awsappstream.Client
}

func (r *associateApplicationEntitlementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_associate_application_entitlement"
}

func (r *associateApplicationEntitlementResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.appstreamClient = meta.appstream
}

func (r *associateApplicationEntitlementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	stackName, entitlementName, applicationIdentifier, err := parseAssociateApplicationEntitlementID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier format: <stack_name>|<entitlement_name>|<application_identifier>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("stack_name"), stackName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("entitlement_name"), entitlementName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_identifier"), applicationIdentifier)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
