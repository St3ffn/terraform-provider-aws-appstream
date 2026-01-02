// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

	if config.AuthenticationType.IsNull() || config.AuthenticationType.IsUnknown() ||
		config.UserName.IsNull() || config.UserName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read user because authentication_type and user_name must be set and known.",
		)
		return
	}

	authenticationType := config.AuthenticationType.ValueString()
	userName := config.UserName.ValueString()

	state, diags := ds.readUser(ctx, authenticationType, userName)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state == nil {
		resp.Diagnostics.AddError(
			"AWS AppStream User Not Found",
			fmt.Sprintf("No user %q was found with authentication type %q.", userName, authenticationType),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (ds *dataSource) readUser(ctx context.Context, authenticationType, userName string) (*dataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var nextToken *string

	for {
		out, err := ds.appstreamClient.DescribeUsers(ctx, &awsappstream.DescribeUsersInput{
			AuthenticationType: awstypes.AuthenticationType(authenticationType),
			NextToken:          nextToken,
		})
		if err != nil {
			if util.IsContextCanceled(err) {
				return nil, diags
			}

			diags.AddError(
				"Error Reading AWS AppStream User",
				fmt.Sprintf(
					"Could not read user %q with authentication type %q: %v",
					userName, authenticationType, err,
				),
			)
			return nil, diags
		}

		for _, user := range out.Users {
			if user.UserName != nil && aws.ToString(user.UserName) == userName {
				awsAuthType := string(user.AuthenticationType)
				awsUserName := aws.ToString(user.UserName)

				state := &dataSourceModel{
					ID:                 types.StringValue(buildID(awsAuthType, awsUserName)),
					AuthenticationType: types.StringValue(awsAuthType),
					UserName:           types.StringValue(awsUserName),
					FirstName:          util.StringOrNull(user.FirstName),
					LastName:           util.StringOrNull(user.LastName),
					Enabled:            util.BoolOrNull(user.Enabled),
					Status:             util.StringOrNull(user.Status),
					ARN:                util.StringOrNull(user.Arn),
					CreatedTime:        util.StringFromTime(user.CreatedTime),
				}
				return state, diags
			}
		}

		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	return nil, diags
}
