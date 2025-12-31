// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
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

	if state.ID.IsNull() || state.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete app block because id must be known.",
		)
		return
	}

	arn := state.ID.ValueString()

	name, err := appBlockNameFromARN(arn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			fmt.Sprintf("Could not parse app block name from ARN %q: %v", arn, err),
		)
		return
	}

	_, err = r.appstreamClient.DeleteAppBlock(ctx, &awsappstream.DeleteAppBlockInput{
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
			"Error Deleting AWS AppStream App Block",
			fmt.Sprintf("Could not delete app block %q: %v", name, err),
		)
		return
	}
}

func appBlockNameFromARN(arn string) (string, error) {
	// expected arn:aws:appstream:<region>:<account>:app-block/<name>
	const prefix = "app-block/"
	idx := strings.LastIndex(arn, prefix)
	if idx == -1 || idx+len(prefix) >= len(arn) {
		return "", fmt.Errorf("invalid app block ARN format")
	}

	return arn[idx+len(prefix):], nil
}
