// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder

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
	_ tfresource.ResourceWithModifyPlan     = &resource{}
	_ tfresource.ResourceWithValidateConfig = &resource{}
	_ tfresource.ResourceWithImportState    = &resource{}
)

// NewResource constructs the resource implementation used by Terraform for this type.
func NewResource() tfresource.Resource {
	return &resource{}
}

type resource struct {
	appstreamClient *awsappstream.Client
	tags            *tags.TagManager
}

// ValidateConfig enforces image source selection rules for image builders.
// Exactly one of image_name or image_arn must be set.
func (r *resource) ValidateConfig(ctx context.Context, req tfresource.ValidateConfigRequest, resp *tfresource.ValidateConfigResponse) {
	var config resourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasImageName := !config.ImageName.IsNull()
	hasImageArn := !config.ImageARN.IsNull()

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
}

// Metadata registers this component's Terraform type name.
// Terraform uses it to bind configuration blocks to this implementation.
func (r *resource) Metadata(_ context.Context, req tfresource.MetadataRequest, resp *tfresource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image_builder"
}

// Configure reads provider metadata, validates the expected metadata type and required clients,
// and stores them on the receiver for subsequent operations.
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

// ImportState expects <image_builder_name> and seeds both name and id.
func (r *resource) ImportState(ctx context.Context, req tfresource.ImportStateRequest, resp *tfresource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier format: <image_builder_name>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
