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

var (
	_ resource.Resource                = &associateFleetStackResource{}
	_ resource.ResourceWithConfigure   = &associateFleetStackResource{}
	_ resource.ResourceWithImportState = &associateFleetStackResource{}
)

func NewAssociateFleetStackResource() resource.Resource {
	return &associateFleetStackResource{}
}

type associateFleetStackResource struct {
	appstreamClient *awsappstream.Client
}

func (r *associateFleetStackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_associate_fleet_stack"
}

func (r *associateFleetStackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *associateFleetStackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fleetName, stackName, err := parseAssociateFleetStackID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier format: <fleet_name>|<stack_name>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("fleet_name"), fleetName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("stack_name"), stackName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
