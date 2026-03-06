// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Create imports an EC2 AMI into AppStream, applies tags, waits for AVAILABLE state,
// and then reads back the remote image to persist authoritative Terraform state.
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
		plan.IAMRoleARN.IsNull() || plan.IAMRoleARN.IsUnknown() ||
		plan.SourceAmiID.IsNull() || plan.SourceAmiID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create imported image because name, iam_role_arn, and source_ami_id must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	input := &awsappstream.CreateImportedImageInput{
		Name:        util.StringPointerOrNil(plan.Name),
		IamRoleArn:  util.StringPointerOrNil(plan.IAMRoleARN),
		SourceAmiId: util.StringPointerOrNil(plan.SourceAmiID),
	}

	input.DisplayName = util.StringPointerOrNil(plan.DisplayName)
	input.Description = util.StringPointerOrNil(plan.Description)

	if !plan.AgentSoftwareVersion.IsNull() && !plan.AgentSoftwareVersion.IsUnknown() {
		input.AgentSoftwareVersion = awstypes.AgentSoftwareVersion(plan.AgentSoftwareVersion.ValueString())
	}

	if !plan.RuntimeValidationConfig.IsNull() && !plan.RuntimeValidationConfig.IsUnknown() {
		input.RuntimeValidationConfig = expandRuntimeValidationConfig(ctx, plan.RuntimeValidationConfig, &resp.Diagnostics)
	}

	if !plan.AppCatalogConfig.IsNull() && !plan.AppCatalogConfig.IsUnknown() {
		input.AppCatalogConfig = expandAppCatalogConfig(ctx, plan.AppCatalogConfig, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	out, err := util.RetryOnValue(
		ctx,
		func(ctx context.Context) (*awsappstream.CreateImportedImageOutput, error) {
			return r.appstreamClient.CreateImportedImage(ctx, input)
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateImportedImage.html
		util.WithRetryOnFns(
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
			util.IsInvalidRoleException,
		),
	)
	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Imported Image Already Exists",
				fmt.Sprintf(
					"An imported image named %q already exists. To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					name, name,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Imported Image",
			fmt.Sprintf("Could not create imported image %q: %v", name, err),
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

	err = r.ensureAvailable(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Imported Image",
			fmt.Sprintf("Imported image %q did not reach AVAILABLE state in time: %v", name, err),
		)
		return
	}

	newState, diags := r.readImportedImage(ctx, plan)
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
				return fmt.Errorf("%w: current=NOT_FOUND", errUnexpectedImportedImageState)
			}

			image := out.Images[0]
			switch image.State {
			case awstypes.ImageStateAvailable:
				return nil

			case awstypes.ImageStateFailed:
				return fmt.Errorf("current=%s", image.State)

			default:
				return fmt.Errorf("%w: current=%s", errUnexpectedImportedImageState, image.State)
			}
		},
		util.WithTimeout(waitAvailableRetryTimeout),
		util.WithInitBackoff(waitAvailableRetryInitBackoff),
		util.WithMaxBackoff(waitAvailableRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeImages.html
		util.WithRetryOnFns(
			isUnexpectedImportedImageStateError,
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)
}
