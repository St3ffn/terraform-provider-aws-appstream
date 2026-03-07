// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Create starts an AppStream image copy from source_image_name into destination_region,
// waits until the destination image reaches AVAILABLE, applies tags to the copied image ARN,
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

	if plan.Name.IsNull() || plan.Name.IsUnknown() ||
		plan.DestinationRegion.IsNull() || plan.DestinationRegion.IsUnknown() ||
		plan.SourceImageName.IsNull() || plan.SourceImageName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create copied image because name, destination_region, and source_image_name must be known.",
		)
		return
	}

	name := plan.Name.ValueString()
	destinationRegion := plan.DestinationRegion.ValueString()
	sourceImageName := plan.SourceImageName.ValueString()

	input := &awsappstream.CopyImageInput{
		DestinationImageName:        aws.String(name),
		DestinationRegion:           aws.String(destinationRegion),
		SourceImageName:             aws.String(sourceImageName),
		DestinationImageDescription: util.StringPointerOrNil(plan.Description),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.CopyImage(ctx, input)
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CopyImage.html
		util.WithRetryOnFns(
			util.IsResourceNotFoundException,
		),
	)
	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Copied Image Already Exists",
				fmt.Sprintf(
					"A copied image named %q already exists in %q. To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					name, destinationRegion, buildID(name, destinationRegion),
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Copied Image",
			fmt.Sprintf("Could not create copied image %q in %q: %v", name, destinationRegion, err),
		)
		return
	}

	imageARN, err := r.ensureAvailable(ctx, name, destinationRegion)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Copied Image",
			fmt.Sprintf(
				"Copied image %q in %q did not reach AVAILABLE state in time: %v",
				name, destinationRegion, err,
			),
		)
		return
	}
	if imageARN == nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Copied Image",
			fmt.Sprintf(
				"Could not determine ARN for copied image %q in %q after it reached AVAILABLE state. "+
					"The image cannot be tagged without an ARN.",
				name, destinationRegion,
			),
		)
		return
	}

	_, tagDiags := r.tags.Apply(
		ctx,
		aws.ToString(imageARN),
		plan.Tags,
		func(o *awstaggingapi.Options) {
			o.Region = destinationRegion
		},
	)
	resp.Diagnostics.Append(tagDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readCopiedImage(ctx, plan)
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

func (r *resource) ensureAvailable(ctx context.Context, name, destinationRegion string) (imageARN *string, err error) {
	return util.RetryOnValue(
		ctx,
		func(ctx context.Context) (*string, error) {
			out, err := r.appstreamClient.DescribeImages(
				ctx,
				&awsappstream.DescribeImagesInput{
					Names: []string{name},
				},
				func(o *awsappstream.Options) {
					o.Region = destinationRegion
				},
			)
			if err != nil {
				return nil, err
			}

			if len(out.Images) == 0 {
				return nil, fmt.Errorf("%w: current=NOT_FOUND", errUnexpectedCopiedImageState)
			}

			image := out.Images[0]
			switch image.State {
			case awstypes.ImageStateAvailable:
				return image.Arn, nil

			case awstypes.ImageStateFailed:
				return image.Arn, fmt.Errorf("current=%s", image.State)

			default:
				return nil, fmt.Errorf("%w: current=%s", errUnexpectedCopiedImageState, image.State)
			}
		},
		util.WithTimeout(waitAvailableRetryTimeout),
		util.WithInitBackoff(waitAvailableRetryInitBackoff),
		util.WithMaxBackoff(waitAvailableRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeImages.html
		util.WithRetryOnFns(
			isUnexpectedCopiedImageStateError,
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)
}
