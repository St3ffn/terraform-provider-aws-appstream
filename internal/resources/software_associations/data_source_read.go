// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package software_associations

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if config.AssociatedResource.IsNull() || config.AssociatedResource.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read software associations because associated_resource must be set and known.",
		)
		return
	}

	associatedResource := config.AssociatedResource.ValueString()

	var all []awstypes.SoftwareAssociations
	var nextToken *string

	for {
		out, err := ds.appstreamClient.DescribeSoftwareAssociations(ctx, &awsappstream.DescribeSoftwareAssociationsInput{
			AssociatedResource: aws.String(associatedResource),
			NextToken:          nextToken,
		})
		if err != nil {
			if util.IsContextCanceled(err) {
				return
			}

			if util.IsAppStreamNotFound(err) {
				state := &dataSourceModel{
					AssociatedResource:   types.StringValue(associatedResource),
					SoftwareAssociations: types.SetNull(softwareAssociationObjectType),
				}
				resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
				return
			}

			resp.Diagnostics.AddError(
				"Error Reading AWS AppStream Software Associations",
				fmt.Sprintf("Could not read software associations for %q: %v", associatedResource, err),
			)
			return
		}

		if out.AssociatedResource == nil {
			resp.Diagnostics.AddError(
				"Unexpected AWS Response",
				fmt.Sprintf("Software associations for %q was returned without required identifiers.", associatedResource),
			)
			return
		}

		all = append(all, out.SoftwareAssociations...)

		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	state := &dataSourceModel{
		AssociatedResource:   types.StringValue(associatedResource),
		SoftwareAssociations: flattenSoftwareAssociations(ctx, all, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
