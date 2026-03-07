// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/tags"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Read fetches the current remote object and refreshes Terraform state from AWS.
// When the object no longer exists remotely, the resource is removed from state
// to converge Terraform with external deletions.
func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.Name.IsNull() || state.Name.IsUnknown() ||
		state.DestinationRegion.IsNull() || state.DestinationRegion.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attributes name and destination_region are missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	newState, diags := r.readCopiedImage(ctx, state)
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

func (r *resource) readCopiedImage(ctx context.Context, prior resourceModel) (*resourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	name := prior.Name.ValueString()
	destinationRegion := prior.DestinationRegion.ValueString()

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
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
			return nil, diags
		}
		diags.AddError(
			"Error Reading AWS AppStream Copied Image",
			fmt.Sprintf("Could not read copied image %q: %v", name, err),
		)
		return nil, diags
	}

	if len(out.Images) == 0 {
		return nil, diags
	}

	image := out.Images[0]
	if image.Name == nil {
		return nil, diags
	}

	sourceImageName := util.ARNResourceSuffixOrNil(image.BaseImageArn, "appstream", "image/")

	state := &resourceModel{
		ID:                          types.StringValue(buildID(name, destinationRegion)),
		Name:                        types.StringValue(aws.ToString(image.Name)),
		Description:                 util.StringOrNull(image.Description),
		DestinationRegion:           prior.DestinationRegion,
		SourceImageName:             util.FlattenStateOwnedString(prior.SourceImageName, sourceImageName),
		Tags:                        types.MapNull(types.StringType),
		TagsAll:                     types.MapNull(types.StringType),
		ARN:                         util.StringOrNull(image.Arn),
		State:                       types.StringValue(string(image.State)),
		StateChangeReason:           flattenStateChangeReason(ctx, image.StateChangeReason, &diags),
		CreatedTime:                 util.StringFromTime(image.CreatedTime),
		Platform:                    types.StringValue(string(image.Platform)),
		Visibility:                  types.StringValue(string(image.Visibility)),
		BaseImageARN:                util.StringOrNull(image.BaseImageArn),
		ImageBuilderSupported:       util.BoolOrNull(image.ImageBuilderSupported),
		ImageBuilderName:            util.StringOrNull(image.ImageBuilderName),
		Applications:                flattenApplications(ctx, image.Applications, &diags),
		AppstreamAgentVersion:       util.StringOrNull(image.AppstreamAgentVersion),
		PublicBaseImageReleasedDate: util.StringFromTime(image.PublicBaseImageReleasedDate),
		ImagePermissions:            flattenImagePermissions(ctx, image.ImagePermissions, &diags),
		ImageErrors:                 flattenImageErrors(ctx, image.ImageErrors, &diags),
		LatestAppstreamAgentVersion: types.StringValue(string(image.LatestAppstreamAgentVersion)),
		SupportedInstanceFamilies:   util.SetStringOrNull(ctx, image.SupportedInstanceFamilies, &diags),
		DynamicAppProvidersEnabled:  types.StringValue(string(image.DynamicAppProvidersEnabled)),
		ImageSharedWithOthers:       types.StringValue(string(image.ImageSharedWithOthers)),
		ManagedSoftwareIncluded:     util.BoolOrNull(image.ManagedSoftwareIncluded),
		ImageType:                   types.StringValue(string(image.ImageType)),
	}

	if !state.ARN.IsNull() {
		allTags, allTagDiags := r.tags.ReadAll(
			ctx,
			state.ARN.ValueString(),
			func(o *awstaggingapi.Options) {
				o.Region = destinationRegion
			},
		)
		diags.Append(allTagDiags...)
		state.TagsAll = allTags

		resourceTags, resourceTagDiags := tags.ResourceTags(ctx, prior.Tags, allTags)
		diags.Append(resourceTagDiags...)
		state.Tags = resourceTags
	}

	if diags.HasError() {
		return nil, diags
	}
	return state, diags
}
