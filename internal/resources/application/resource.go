// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/path"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/metadata"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/tags"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var (
	_ tfresource.Resource                = &resource{}
	_ tfresource.ResourceWithConfigure   = &resource{}
	_ tfresource.ResourceWithImportState = &resource{}
)

func NewResource() tfresource.Resource {
	return &resource{}
}

type resource struct {
	appstreamClient *awsappstream.Client
	tags            *tags.TagManager
}

func (r *resource) Metadata(_ context.Context, req tfresource.MetadataRequest, resp *tfresource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
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
	if err := util.ValidateARNValue(req.ID, "appstream", "application/"); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected AppStream application ARN: %v", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
