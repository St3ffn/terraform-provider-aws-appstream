// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var stateChangeReasonObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"code":    types.StringType,
		"message": types.StringType,
	},
}

var applicationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":              types.StringType,
		"display_name":      types.StringType,
		"icon_url":          types.StringType,
		"launch_path":       types.StringType,
		"launch_parameters": types.StringType,
		"enabled":           types.BoolType,
		"metadata":          types.MapType{ElemType: types.StringType},
		"working_directory": types.StringType,
		"description":       types.StringType,
		"arn":               types.StringType,
		"app_block_arn":     types.StringType,
		"icon_s3_location":  iconS3LocationObjectType,
		"platforms":         types.SetType{ElemType: types.StringType},
		"instance_families": types.SetType{ElemType: types.StringType},
		"created_time":      types.StringType,
	},
}

var iconS3LocationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"s3_bucket": types.StringType,
		"s3_key":    types.StringType,
	},
}

var imagePermissionsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"allow_fleet":         types.BoolType,
		"allow_image_builder": types.BoolType,
	},
}

var imageErrorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":      types.StringType,
		"error_message":   types.StringType,
		"error_timestamp": types.StringType,
	},
}

func flattenStateChangeReason(
	ctx context.Context, awsReason *awstypes.ImageStateChangeReason, diags *diag.Diagnostics,
) types.Object {

	if awsReason == nil {
		return types.ObjectNull(stateChangeReasonObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		stateChangeReasonObjectType.AttrTypes,
		stateChangeReasonModel{
			Code:    types.StringValue(string(awsReason.Code)),
			Message: util.StringOrNull(awsReason.Message),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(stateChangeReasonObjectType.AttrTypes)
	}

	return obj
}

func flattenApplications(ctx context.Context, awsApps []awstypes.Application, diags *diag.Diagnostics) types.Set {

	if len(awsApps) == 0 {
		return types.SetNull(applicationObjectType)
	}

	out := make([]applicationModel, 0, len(awsApps))
	for _, app := range awsApps {
		out = append(out, applicationModel{
			Name:             util.StringOrNull(app.Name),
			DisplayName:      util.StringOrNull(app.DisplayName),
			IconURL:          util.StringOrNull(app.IconURL),
			LaunchPath:       util.StringOrNull(app.LaunchPath),
			LaunchParameters: util.StringOrNull(app.LaunchParameters),
			Enabled:          util.BoolOrNull(app.Enabled),
			Metadata:         util.MapStringOrNull(ctx, app.Metadata, diags),
			WorkingDirectory: util.StringOrNull(app.WorkingDirectory),
			Description:      util.StringOrNull(app.Description),
			ARN:              util.StringOrNull(app.Arn),
			AppBlockARN:      util.StringOrNull(app.AppBlockArn),
			IconS3Location:   flattenIconS3Location(ctx, app.IconS3Location, diags),
			Platforms:        util.SetEnumStringOrNull(ctx, app.Platforms, diags),
			InstanceFamilies: util.SetStringOrNull(ctx, app.InstanceFamilies, diags),
			CreatedTime:      util.StringFromTime(app.CreatedTime),
		})
	}

	setVal, d := types.SetValueFrom(ctx, applicationObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(applicationObjectType)
	}

	return setVal
}

func flattenIconS3Location(
	ctx context.Context, awsS3Location *awstypes.S3Location, diags *diag.Diagnostics,
) types.Object {

	if awsS3Location == nil {
		return types.ObjectNull(iconS3LocationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		iconS3LocationObjectType.AttrTypes,
		iconS3LocationModel{
			S3Bucket: util.StringOrNull(awsS3Location.S3Bucket),
			S3Key:    util.StringOrNull(awsS3Location.S3Key),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(iconS3LocationObjectType.AttrTypes)
	}

	return obj
}

func flattenImagePermissions(
	ctx context.Context, awsPerms *awstypes.ImagePermissions, diags *diag.Diagnostics,
) types.Object {

	if awsPerms == nil {
		return types.ObjectNull(imagePermissionsObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		imagePermissionsObjectType.AttrTypes,
		imagePermissionsModel{
			AllowFleet:        util.BoolOrNull(awsPerms.AllowFleet),
			AllowImageBuilder: util.BoolOrNull(awsPerms.AllowImageBuilder),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(imagePermissionsObjectType.AttrTypes)
	}

	return obj
}

func flattenImageErrors(ctx context.Context, awsErrors []awstypes.ResourceError, diags *diag.Diagnostics) types.Set {
	if len(awsErrors) == 0 {
		return types.SetNull(imageErrorObjectType)
	}

	out := make([]imageErrorModel, 0, len(awsErrors))
	for _, e := range awsErrors {
		out = append(out, imageErrorModel{
			ErrorCode:      types.StringValue(string(e.ErrorCode)),
			ErrorMessage:   util.StringOrNull(e.ErrorMessage),
			ErrorTimestamp: util.StringFromTime(e.ErrorTimestamp),
		})
	}

	setVal, d := types.SetValueFrom(ctx, imageErrorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(imageErrorObjectType)
	}

	return setVal
}
