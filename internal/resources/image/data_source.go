// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/metadata"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/tags"
)

var (
	_ datasource.DataSource                   = &dataSource{}
	_ datasource.DataSourceWithValidateConfig = &dataSource{}
	_ datasource.DataSourceWithConfigure      = &dataSource{}
)

func NewDataSource() datasource.DataSource {
	return &dataSource{}
}

type dataSource struct {
	appstreamClient *awsappstream.Client
	tags            *tags.TagManager
}

func (ds *dataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var config model

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasArn := !config.ARN.IsNull() && !config.ARN.IsUnknown()
	hasName := !config.Name.IsNull() && !config.Name.IsUnknown()
	hasNameRegex := !config.NameRegex.IsNull() && !config.NameRegex.IsUnknown()

	selectionCount := 0
	if hasArn {
		selectionCount++
	}
	if hasName {
		selectionCount++
	}
	if hasNameRegex {
		selectionCount++
	}

	if selectionCount == 0 || selectionCount > 1 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Exactly one of `arn`, `name`, or `name_regex` must be specified.",
		)
	}
}

func (ds *dataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (ds *dataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*metadata.Metadata)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Metadata, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if meta.Appstream == nil {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *Metadata.Appstream, got: nil. Please report this issue to the provider developers.",
		)
		return
	}

	if meta.Tagging == nil {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *Metadata.Tagging, got: nil. Please report this issue to the provider developers.",
		)
		return
	}

	ds.appstreamClient = meta.Appstream
	ds.tags = tags.NewTagManager(meta.Tagging, meta.DefaultTags)
}
