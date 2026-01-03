// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	if plan.AuthenticationType.IsNull() || plan.AuthenticationType.IsUnknown() ||
		plan.UserName.IsNull() || plan.UserName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create user because authentication_type and user_name must be known.",
		)
		return
	}

	authenticationType := plan.AuthenticationType.ValueString()
	userName := plan.UserName.ValueString()

	input := &awsappstream.CreateUserInput{
		AuthenticationType: awstypes.AuthenticationType(authenticationType),
		UserName:           aws.String(userName),
		FirstName:          util.StringPointerOrNil(plan.FirstName),
		LastName:           util.StringPointerOrNil(plan.LastName),
	}

	if !plan.MessageAction.IsNull() && !plan.MessageAction.IsUnknown() {
		input.MessageAction = awstypes.MessageAction(plan.MessageAction.ValueString())
	}

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.CreateUser(ctx, input)
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateUser.html
		util.WithRetryOnFns(
			util.IsOperationNotPermittedException,
		),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream User Already Exists",
				fmt.Sprintf(
					"A user named %q already exists with authentication type %q. "+
						"To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					userName, authenticationType, buildID(authenticationType, userName),
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream User",
			fmt.Sprintf("Could not create user %q with authentication type %q: %v",
				userName, authenticationType, err,
			),
		)
		return
	}

	// aws creates users enabled by default; apply desired enabled state if false
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() &&
		!plan.Enabled.ValueBool() {

		err = util.RetryOn(
			ctx,
			func(ctx context.Context) error {
				var err error
				_, err = r.appstreamClient.DisableUser(ctx, &awsappstream.DisableUserInput{
					AuthenticationType: awstypes.AuthenticationType(authenticationType),
					UserName:           aws.String(userName),
				})
				return err
			},
			util.WithTimeout(disableRetryTimeout),
			util.WithInitBackoff(disableRetryInitBackoff),
			util.WithMaxBackoff(disableRetryMaxBackoff),
			// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DisableUser.html
			util.WithRetryOnFns(
				util.IsResourceNotFoundException,
			),
		)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating AWS AppStream User",
				fmt.Sprintf("Could not disable user %q with authentication type %q: %v",
					userName, authenticationType, err,
				),
			)
			return
		}
	}

	newState, diags := r.readUser(ctx, plan)
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
