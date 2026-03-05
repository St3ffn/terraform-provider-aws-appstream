// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Delete calls the AWS delete/disassociate API and treats not-found
// responses as already deleted so destroy remains idempotent.
// Terraform state is then cleared by the framework lifecycle.
func (r *resource) Delete(ctx context.Context, req tfresource.DeleteRequest, resp *tfresource.DeleteResponse) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete application because id must be known.",
		)
		return
	}

	arn := state.ID.ValueString()

	name, err := applicationNameFromARN(arn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			fmt.Sprintf("Could not parse application name from ARN %q: %v", arn, err),
		)
		return
	}

	_, err = r.appstreamClient.DeleteApplication(ctx, &awsappstream.DeleteApplicationInput{
		Name: aws.String(name),
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
			"Error Deleting AWS AppStream Application",
			fmt.Sprintf("Could not delete application %q: %v", name, err),
		)
		return
	}
}

func applicationNameFromARN(arn string) (string, error) {
	// expected arn:aws:appstream:<region>:<account>:application/<name>
	const prefix = "application/"
	idx := strings.LastIndex(arn, prefix)
	if idx == -1 || idx+len(prefix) >= len(arn) {
		return "", fmt.Errorf("invalid application ARN format")
	}

	return arn[idx+len(prefix):], nil
}
