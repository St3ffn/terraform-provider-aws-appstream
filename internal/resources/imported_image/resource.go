// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

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
	_ tfresource.Resource                = &resource{}
	_ tfresource.ResourceWithConfigure   = &resource{}
	_ tfresource.ResourceWithModifyPlan  = &resource{}
	_ tfresource.ResourceWithImportState = &resource{}
)

// NewResource constructs the resource implementation used by Terraform for this type.
func NewResource() tfresource.Resource {
	return &resource{}
}

type resource struct {
	appstreamClient *awsappstream.Client
	tags            *tags.TagManager
}

// Metadata registers this component's Terraform type name.
// Terraform uses it to bind configuration blocks to this implementation.
func (r *resource) Metadata(_ context.Context, req tfresource.MetadataRequest, resp *tfresource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_imported_image"
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

// ImportState expects <name>|<iam_role_arn>|<source_ami_id>
// and seeds name, iam_role_arn, source_ami_id, and id from that identifier.
func (r *resource) ImportState(ctx context.Context, req tfresource.ImportStateRequest, resp *tfresource.ImportStateResponse) {
	name, iamRoleARN, sourceAmiID, err := parseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier format: <name>|<iam_role_arn>|<source_ami_id>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("iam_role_arn"), iamRoleARN)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_ami_id"), sourceAmiID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
