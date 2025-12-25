// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &applicationDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationDataSource{}
)

func NewApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

type applicationDataSource struct {
	appstreamClient *awsappstream.Client
	tags            *tagManager
}

func (ds *applicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (ds *applicationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*metadata)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *metadata, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if meta.appstream == nil {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *metadata.appstream, got: nil. Please report this issue to the provider developers.",
		)
		return
	}

	if meta.tagging == nil {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *metadata.tagging, got: nil. Please report this issue to the provider developers.",
		)
		return
	}

	ds.appstreamClient = meta.appstream
	ds.tags = newTagManager(meta.tagging, meta.defaultTags)
}
