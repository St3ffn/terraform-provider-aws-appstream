// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import (
	"context"
	"errors"
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
			"Cannot delete image builder because name must be known.",
		)
		return
	}

	name := state.Name.ValueString()

	err := r.deleteImageBuilder(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Image Builder",
			fmt.Sprintf("Could not delete image builder %q: %v", name, err),
		)
		return
	}
}

var ErrUnexpectedImageBuilderState = errors.New("unexpected image builder state")

func (r *resource) deleteImageBuilder(ctx context.Context, name string) error {
	return util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			out, err := r.appstreamClient.DescribeImageBuilders(ctx, &awsappstream.DescribeImageBuildersInput{
				Names: []string{name},
			})
			if err != nil {
				if util.IsAppStreamNotFound(err) {
					// already deleted
					return nil
				}
				return err
			}

			if len(out.ImageBuilders) == 0 {
				return nil
			}

			state := out.ImageBuilders[0].State

			switch state {
			case awstypes.ImageBuilderStateRunning:
				// stoppable state
				_, err = r.appstreamClient.StopImageBuilder(ctx, &awsappstream.StopImageBuilderInput{
					Name: aws.String(name),
				})
				if err != nil {
					if util.IsAppStreamNotFound(err) {
						// already deleted
						return nil
					}
					return err
				}
				// retry as we just stopped the image builder
				return fmt.Errorf("%w: current=%s", ErrUnexpectedImageBuilderState, state)

			case awstypes.ImageBuilderStateStopped, awstypes.ImageBuilderStateFailed:
				// deletable state
				_, err = r.appstreamClient.DeleteImageBuilder(ctx, &awsappstream.DeleteImageBuilderInput{
					Name: aws.String(name),
				})
				if err != nil {
					if util.IsAppStreamNotFound(err) {
						// already deleted
						return nil
					}
					return err
				}
				// wait until resource is in state deleting or gone
				return fmt.Errorf("%w: current=%s", ErrUnexpectedImageBuilderState, state)
			default:
				return fmt.Errorf("%w: current=%s", ErrUnexpectedImageBuilderState, state)
			}
		},
		util.WithTimeout(imageBuilderWaitTimeout),
		util.WithInitBackoff(imageBuilderWaitInitBackoff),
		util.WithMaxBackoff(imageBuilderWaitMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeImageBuilders.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_StopImageBuilder.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DeleteImageBuilder.html
		util.WithRetryOnFns(
			func(err error) bool {
				return errors.Is(err, ErrUnexpectedImageBuilderState)
			},
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
		),
	)
}
