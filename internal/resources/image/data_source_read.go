// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config model

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	images, err := ds.listImages(ctx, &config)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Image Not Found",
				"No image matched the given criteria.",
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Images",
			fmt.Sprintf("Could not read images: %v", err),
		)
		return
	}

	if len(images) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream Image Not Found",
			"No image matched the given criteria.",
		)
		return
	}

	selected, err := selectSingleImage(images, config.MostRecent)
	if err != nil {
		resp.Diagnostics.AddError(
			"Multiple AppStream Images Found",
			err.Error(),
		)
		return
	}

	state := &model{
		ID:                          types.StringValue(aws.ToString(selected.Arn)),
		ARN:                         types.StringValue(aws.ToString(selected.Arn)),
		Name:                        types.StringValue(aws.ToString(selected.Name)),
		Visibility:                  types.StringValue(string(selected.Visibility)),
		BaseImageARN:                util.StringOrNull(selected.BaseImageArn),
		DisplayName:                 util.StringOrNull(selected.DisplayName),
		State:                       types.StringValue(string(selected.State)),
		ImageBuilderSupported:       util.BoolOrNull(selected.ImageBuilderSupported),
		ImageBuilderName:            util.StringOrNull(selected.ImageBuilderName),
		Platform:                    types.StringValue(string(selected.Platform)),
		Description:                 util.StringOrNull(selected.Description),
		StateChangeReason:           flattenStateChangeReason(ctx, selected.StateChangeReason, &resp.Diagnostics),
		Applications:                flattenApplications(ctx, selected.Applications, &resp.Diagnostics),
		CreatedTime:                 util.StringFromTime(selected.CreatedTime),
		PublicBaseImageReleasedDate: util.StringFromTime(selected.PublicBaseImageReleasedDate),
		AppstreamAgentVersion:       util.StringOrNull(selected.AppstreamAgentVersion),
		ImagePermissions:            flattenImagePermissions(ctx, selected.ImagePermissions, &resp.Diagnostics),
		ImageErrors:                 flattenImageErrors(ctx, selected.ImageErrors, &resp.Diagnostics),
		LatestAppstreamAgentVersion: types.StringValue(string(selected.LatestAppstreamAgentVersion)),
		SupportedInstanceFamilies:   util.SetStringOrNull(ctx, selected.SupportedInstanceFamilies, &resp.Diagnostics),
		DynamicAppProvidersEnabled:  types.StringValue(string(selected.DynamicAppProvidersEnabled)),
		ImageSharedWithOthers:       types.StringValue(string(selected.ImageSharedWithOthers)),
		ManagedSoftwareIncluded:     util.BoolOrNull(selected.ManagedSoftwareIncluded),
		ImageType:                   types.StringValue(string(selected.ImageType)),
		Tags:                        types.MapNull(types.StringType),
	}

	if !state.ARN.IsNull() && selected.Visibility != awstypes.VisibilityTypePublic {
		tags, diags := ds.tags.Read(ctx, state.ARN.ValueString())
		resp.Diagnostics.Append(diags...)
		state.Tags = tags
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (ds *dataSource) listImages(ctx context.Context, config *model) ([]awstypes.Image, error) {
	var (
		arns  []string
		names []string
	)

	if !config.ARN.IsNull() && !config.ARN.IsUnknown() {
		arns = []string{config.ARN.ValueString()}
	}

	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		names = []string{config.Name.ValueString()}
	}

	var regex *regexp.Regexp
	if !config.NameRegex.IsNull() && !config.NameRegex.IsUnknown() {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			return nil, err
		}
		regex = r
	}

	var out []awstypes.Image
	var nextToken *string

	for {
		input := &awsappstream.DescribeImagesInput{
			Arns:      arns,
			Names:     names,
			NextToken: nextToken,
		}
		if !config.Visibility.IsNull() && !config.Visibility.IsUnknown() {
			input.Type = awstypes.VisibilityType(config.Visibility.ValueString())
		}

		resp, err := ds.appstreamClient.DescribeImages(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, image := range resp.Images {
			if regex != nil {
				if image.Name == nil || !regex.MatchString(*image.Name) {
					continue
				}
			}
			out = append(out, image)
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return out, nil
}

func selectSingleImage(images []awstypes.Image, mostRecent types.Bool) (*awstypes.Image, error) {
	if len(images) == 1 {
		return &images[0], nil
	}

	useMostRecent := !mostRecent.IsNull() && !mostRecent.IsUnknown() && mostRecent.ValueBool()

	if !useMostRecent {
		return nil, fmt.Errorf(
			"multiple images matched the selection criteria; set most_recent = true to select the newest image",
		)
	}

	sort.Slice(images, func(i, j int) bool {
		ti := images[i].CreatedTime
		tj := images[j].CreatedTime

		if ti == nil {
			return false
		}
		if tj == nil {
			return true
		}
		return ti.After(*tj)
	})

	return &images[0], nil
}
