// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"
	"fmt"

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

	if state.DirectoryName.IsNull() || state.DirectoryName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete directory config because directory_name must be known.",
		)
		return
	}

	name := state.DirectoryName.ValueString()

	_, err := r.appstreamClient.DeleteDirectoryConfig(ctx, &awsappstream.DeleteDirectoryConfigInput{
		DirectoryName: aws.String(name),
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
			"Error Deleting AWS AppStream Directory Config",
			fmt.Sprintf("Could not delete directory config %q: %v", name, err),
		)
		return
	}
}
