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

func (r *resource) Delete(ctx context.Context, req tfresource.DeleteRequest, resp *tfresource.DeleteResponse) {
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
			"Cannot delete user because authentication_type and user_name must be known.",
		)
		return
	}

	authenticationType := state.AuthenticationType.ValueString()
	userName := state.UserName.ValueString()

	_, err := r.appstreamClient.DeleteUser(ctx, &awsappstream.DeleteUserInput{
		AuthenticationType: awstypes.AuthenticationType(authenticationType),
		UserName:           aws.String(userName),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		// if it's already gone, that's fine for delete.
		if util.IsAppStreamNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream User",
			fmt.Sprintf("Could not delete user %q with authentication type %q: %v",
				userName, authenticationType, err,
			),
		)
		return
	}
}
