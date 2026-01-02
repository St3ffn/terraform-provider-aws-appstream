// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.AuthenticationType.IsNull() || state.AuthenticationType.IsUnknown() ||
		state.UserName.IsNull() || state.UserName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attributes authentication_type and user_name are missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	newState, diags := r.readUser(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if ctx.Err() != nil {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *resource) readUser(ctx context.Context, prior resourceModel) (*resourceModel, diag.Diagnostics) {
	var state *resourceModel
	var diags diag.Diagnostics

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			var err error
			state, err = r.readUserOnce(ctx, prior)
			return err
		},
		util.WithTimeout(5*time.Minute),
		util.WithInitBackoff(1*time.Second),
		util.WithMaxBackoff(10*time.Second),
		util.WithRetryOnFns(
			isUserNotYetVisibleError,
			util.IsResourceNotFoundException,
		),
	)

	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		authenticationType := prior.AuthenticationType.ValueString()
		userName := prior.UserName.ValueString()

		diags.AddError(
			"Error Reading AWS AppStream User",
			fmt.Sprintf(
				"Could not read user %q with authentication type %q: %v",
				userName, authenticationType, err,
			),
		)
		return nil, diags
	}

	return state, diags
}

func (r *resource) readUserOnce(ctx context.Context, prior resourceModel) (*resourceModel, error) {
	var nextToken *string

	authenticationType := prior.AuthenticationType.ValueString()
	userName := prior.UserName.ValueString()

	for {
		out, err := r.appstreamClient.DescribeUsers(ctx, &awsappstream.DescribeUsersInput{
			AuthenticationType: awstypes.AuthenticationType(authenticationType),
			NextToken:          nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, user := range out.Users {
			if user.UserName != nil && aws.ToString(user.UserName) == userName {
				awsAuthType := string(user.AuthenticationType)
				awsUserName := aws.ToString(user.UserName)

				state := &resourceModel{
					ID:                 types.StringValue(buildID(awsAuthType, awsUserName)),
					AuthenticationType: types.StringValue(awsAuthType),
					UserName:           types.StringValue(awsUserName),
					FirstName:          util.FlattenOwnedString(prior.FirstName, user.FirstName),
					LastName:           util.FlattenOwnedString(prior.LastName, user.LastName),
					MessageAction:      prior.MessageAction,
					Enabled:            util.BoolOrNull(user.Enabled),
					Status:             util.StringOrNull(user.Status),
					ARN:                util.StringOrNull(user.Arn),
					CreatedTime:        util.StringFromTime(user.CreatedTime),
				}
				return state, nil
			}
		}

		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	return nil, errUserNotYetVisible
}
