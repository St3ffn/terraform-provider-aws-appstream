// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_permission

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

	if config.Name.IsNull() || config.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read image permissions because name must be set and known.",
		)
		return
	}

	name := config.Name.ValueString()

	var all []awstypes.SharedImagePermissions
	var nextToken *string

	for {
		out, err := ds.appstreamClient.DescribeImagePermissions(ctx, &awsappstream.DescribeImagePermissionsInput{
			Name:       aws.String(name),
			NextToken:  nextToken,
			MaxResults: aws.Int32(AppStreamMaxResults),
		})
		if err != nil {
			if util.IsContextCanceled(err) {
				return
			}

			if util.IsAppStreamNotFound(err) {
				state := &dataSourceModel{
					Name:        types.StringValue(name),
					Permissions: types.SetNull(imagePermissionEntryObjectType),
				}
				resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
				return
			}

			resp.Diagnostics.AddError(
				"Error Reading AWS AppStream Image Permissions",
				fmt.Sprintf("Could not read image permissions for image %q: %v", name, err),
			)
			return
		}

		if out.Name == nil {
			resp.Diagnostics.AddError(
				"Unexpected AWS Response",
				fmt.Sprintf("Image permission for image %q was returned without required identifiers.", name),
			)
			return
		}

		all = append(all, out.SharedImagePermissionsList...)

		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	state := &dataSourceModel{
		Name:        types.StringValue(name),
		Permissions: flattenImagePermissionEntriesData(ctx, all, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
