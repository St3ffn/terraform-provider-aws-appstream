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

func (r *resource) Update(ctx context.Context, req tfresource.UpdateRequest, resp *tfresource.UpdateResponse) {
	var plan resourceModel
	var state resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
			"Cannot update user because authentication_type and user_name must be known.",
		)
		return
	}

	authenticationType := plan.AuthenticationType.ValueString()
	userName := plan.UserName.ValueString()

	// guard against unexpected identity drift
	if !state.AuthenticationType.IsNull() && !state.AuthenticationType.IsUnknown() &&
		!state.UserName.IsNull() && !state.UserName.IsUnknown() {
		if state.AuthenticationType.ValueString() != authenticationType || state.UserName.ValueString() != userName {
			resp.Diagnostics.AddError(
				"Unexpected Update Request",
				"User identity (authentication_type/user_name) changed during update. This should trigger replacement. Please report this issue.",
			)
			return
		}
	}

	var err error

	if plan.Enabled.IsNull() || plan.Enabled.IsUnknown() {
		// user does not manage enabled
		goto READ
	} else if state.Enabled.IsUnknown() || state.Enabled.IsNull() {
		// first time management
	} else if state.Enabled.ValueBool() == plan.Enabled.ValueBool() {
		// no changes
		goto READ
	}

	if plan.Enabled.ValueBool() {
		_, err = r.appstreamClient.EnableUser(ctx, &awsappstream.EnableUserInput{
			AuthenticationType: awstypes.AuthenticationType(authenticationType),
			UserName:           aws.String(userName),
		})
	} else {
		_, err = r.appstreamClient.DisableUser(ctx, &awsappstream.DisableUserInput{
			AuthenticationType: awstypes.AuthenticationType(authenticationType),
			UserName:           aws.String(userName),
		})
	}

	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		// disappeared, treat as gone and next plan/apply will recreate
		if util.IsAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream User",
			fmt.Sprintf("Could not update user %q with authentication type %q: %v", userName, authenticationType, err),
		)
		return
	}

READ:
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
