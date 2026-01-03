// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

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
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addDiagnostics(state, &resp.Diagnostics, diagnosticDelete)
	if resp.Diagnostics.HasError() {
		return
	}

	stackName := state.StackName.ValueString()
	userName := state.UserName.ValueString()
	authenticationType := state.AuthenticationType.ValueString()

	out, err := r.appstreamClient.BatchDisassociateUserStack(ctx, &awsappstream.BatchDisassociateUserStackInput{
		UserStackAssociations: []awstypes.UserStackAssociation{
			{
				AuthenticationType: awstypes.AuthenticationType(authenticationType),
				StackName:          aws.String(stackName),
				UserName:           aws.String(userName),
			},
		},
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
			"Error Deleting AWS AppStream User Stack Association",
			fmt.Sprintf("Could not disassociate user %q from stack %q: %v", userName, stackName, err),
		)
		return
	}

	for _, e := range out.Errors {
		if e.ErrorCode == awstypes.UserStackAssociationErrorCodeStackNotFound ||
			e.ErrorCode == awstypes.UserStackAssociationErrorCodeUserNameNotFound {
			// already gone, that's fine for delete.
			continue
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream User Stack Association",
			fmt.Sprintf(
				"Could not disassociate user %q from stack %q: %s",
				userName, stackName, aws.ToString(e.ErrorMessage),
			),
		)
		return
	}
}
