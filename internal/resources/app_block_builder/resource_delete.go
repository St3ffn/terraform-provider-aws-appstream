// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

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

	if state.Name.IsNull() || state.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete app block builder because name must be known.",
		)
		return
	}

	name := state.Name.ValueString()

	err := r.deleteAppBlockBuilder(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream App Block Builder",
			fmt.Sprintf("Could not delete app block builder %q: %v", name, err),
		)
		return
	}
}

func (r *resource) deleteAppBlockBuilder(ctx context.Context, name string) error {
	return util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			out, err := r.appstreamClient.DescribeAppBlockBuilders(ctx, &awsappstream.DescribeAppBlockBuildersInput{
				Names: []string{name},
			})
			if err != nil {
				if util.IsAppStreamNotFound(err) {
					// already deleted
					return nil
				}
				return err
			}

			if len(out.AppBlockBuilders) == 0 {
				return nil
			}

			state := out.AppBlockBuilders[0].State

			switch state {
			case awstypes.AppBlockBuilderStateRunning:
				// stoppable state
				_, err = r.appstreamClient.StopAppBlockBuilder(ctx, &awsappstream.StopAppBlockBuilderInput{
					Name: aws.String(name),
				})
				if err != nil {
					if util.IsAppStreamNotFound(err) {
						// already deleted
						return nil
					}
					return err
				}
				// retry as we just stopped the app block builder
				return fmt.Errorf("%w: current=%s", errUnexpectedAppBlockBuilderState, state)

			case awstypes.AppBlockBuilderStateStopped:
				// deletable state
				_, err = r.appstreamClient.DeleteAppBlockBuilder(ctx, &awsappstream.DeleteAppBlockBuilderInput{
					Name: aws.String(name),
				})
				if err != nil {
					if util.IsAppStreamNotFound(err) {
						// already deleted
						return nil
					}
					return err
				}
				// wait until resource is gone
				return fmt.Errorf("%w: current=%s", errUnexpectedAppBlockBuilderState, state)

			default:
				return fmt.Errorf("%w: current=%s", errUnexpectedAppBlockBuilderState, state)
			}
		},
		util.WithTimeout(deleteRetryTimeout),
		util.WithInitBackoff(deleteRetryInitBackoff),
		util.WithMaxBackoff(deleteRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeAppBlockBuilders.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_StopAppBlockBuilder.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DeleteAppBlockBuilder.html
		util.WithRetryOnFns(
			isUnexpectedAppBlockBuilderStateError,
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceInUseException,
		),
	)
}
