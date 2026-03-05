// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package sessions

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/metadata"
)

var (
	_ datasource.DataSource                   = &dataSource{}
	_ datasource.DataSourceWithConfigure      = &dataSource{}
	_ datasource.DataSourceWithValidateConfig = &dataSource{}
)

// NewDataSource constructs the data source implementation used by Terraform for this type.
func NewDataSource() datasource.DataSource {
	return &dataSource{}
}

type dataSource struct {
	appstreamClient *awsappstream.Client
}

// ValidateConfig validates session filter dependencies.
// authentication_type is required whenever user_id is specified.
func (ds *dataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var config dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.UserID.IsNull() && !config.UserID.IsUnknown() {
		if config.AuthenticationType.IsNull() || config.AuthenticationType.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_type"),
				"Invalid Session Filter Configuration",
				"`authentication_type` must be specified when `user_id` is set.",
			)
		}
	}
}

// Metadata registers this component's Terraform type name.
// Terraform uses it to bind configuration blocks to this implementation.
func (ds *dataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sessions"
}

// Configure reads provider metadata, validates the expected metadata type and required clients,
// and stores them on the receiver for subsequent operations.
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

	ds.appstreamClient = meta.Appstream
}
