// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) waitForImageBuilderAssociable(ctx context.Context, name string) error {
	return util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			out, err := r.appstreamClient.DescribeImageBuilders(ctx, &awsappstream.DescribeImageBuildersInput{
				Names: []string{name},
			})
			if err != nil {
				return err
			}

			if len(out.ImageBuilders) == 0 {
				return fmt.Errorf("image builder %q not found", name)
			}

			state := out.ImageBuilders[0].State

			switch state {
			case awstypes.ImageBuilderStateStopped, awstypes.ImageBuilderStateRunning:
				return nil
			case awstypes.ImageBuilderStateFailed, awstypes.ImageBuilderStateDeleting:
				return fmt.Errorf("image builder is in non-associable state %q", state)
			default:
				return fmt.Errorf("%w: current=%s", errUnexpectedImageBuilderState, state)
			}
		},
		util.WithTimeout(waitAssociableRetryTimeout),
		util.WithInitBackoff(waitAssociableRetryInitBackoff),
		util.WithMaxBackoff(waitAssociableRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeImageBuilders.html
		util.WithRetryOnFns(
			isUnexpectedImageBuilderStateError,
		),
	)
}
