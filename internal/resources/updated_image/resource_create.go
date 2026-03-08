// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package updated_image

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Create starts an AppStream image update from existing_image_name to new_image_name,
// waits until the destination image reaches AVAILABLE, applies tags to the updated image ARN,
// and then performs a final read so Terraform state reflects remote values.
func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	if plan.ExistingImageName.IsNull() || plan.ExistingImageName.IsUnknown() ||
		plan.NewImageName.IsNull() || plan.NewImageName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create updated image because existing_image_name and new_image_name must be known.",
		)
		return
	}

	existingImageName := plan.ExistingImageName.ValueString()
	newImageName := plan.NewImageName.ValueString()

	input := &awsappstream.CreateUpdatedImageInput{
		ExistingImageName:   aws.String(existingImageName),
		NewImageName:        aws.String(newImageName),
		NewImageDescription: util.StringPointerOrNil(plan.NewImageDescription),
		NewImageDisplayName: util.StringPointerOrNil(plan.NewImageDisplayName),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	out, err := util.RetryOnValue(
		ctx,
		func(ctx context.Context) (*awsappstream.CreateUpdatedImageOutput, error) {
			return r.appstreamClient.CreateUpdatedImage(ctx, input)
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateUpdatedImage.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
		),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Updated Image Already Exists",
				fmt.Sprintf(
					"A updated image named %q already exists. To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					newImageName, buildID(existingImageName, newImageName),
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Updated Image",
			fmt.Sprintf("Could not create updated image %q from %q: %v", newImageName, existingImageName, err),
		)
		return
	}

	if out.Image != nil && out.Image.Arn != nil {
		_, tagDiags := r.tags.Apply(ctx, aws.ToString(out.Image.Arn), plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	err = r.ensureAvailable(ctx, newImageName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Updated Image",
			fmt.Sprintf("Updated image %q did not reach AVAILABLE state in time: %v", newImageName, err),
		)
		return
	}

	newState, diags := r.readUpdatedImage(ctx, plan)
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

func (r *resource) ensureAvailable(ctx context.Context, name string) error {
	return util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			out, err := r.appstreamClient.DescribeImages(ctx, &awsappstream.DescribeImagesInput{
				Names: []string{name},
			})
			if err != nil {
				return err
			}

			if len(out.Images) == 0 {
				return fmt.Errorf("%w: current=NOT_FOUND", errUnexpectedUpdatedImageState)
			}

			image := out.Images[0]
			switch image.State {
			case awstypes.ImageStateAvailable:
				return nil

			case awstypes.ImageStateFailed:
				return fmt.Errorf("current=%s", image.State)

			default:
				return fmt.Errorf("%w: current=%s", errUnexpectedUpdatedImageState, image.State)
			}
		},
		util.WithTimeout(waitAvailableRetryTimeout),
		util.WithInitBackoff(waitAvailableRetryInitBackoff),
		util.WithMaxBackoff(waitAvailableRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeImages.html
		util.WithRetryOnFns(
			isUnexpectedUpdatedImageStateError,
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)
}
