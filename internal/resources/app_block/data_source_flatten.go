// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func flattenSourceS3LocationData(
	ctx context.Context, awsS3Location *awstypes.S3Location, diags *diag.Diagnostics,
) types.Object {

	if awsS3Location == nil {
		return types.ObjectNull(sourceS3LocationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		sourceS3LocationObjectType.AttrTypes,
		sourceS3LocationModel{
			S3Bucket: util.StringOrNull(awsS3Location.S3Bucket),
			S3Key:    util.StringOrNull(awsS3Location.S3Key),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(sourceS3LocationObjectType.AttrTypes)
	}

	return obj
}

func flattenScriptDetailsData(
	ctx context.Context, awsScriptDetails *awstypes.ScriptDetails, diags *diag.Diagnostics,
) types.Object {

	if awsScriptDetails == nil {
		return types.ObjectNull(scriptDetailsObjectType.AttrTypes)
	}

	detailsModel := scriptDetailsModel{
		ExecutablePath:       util.StringOrNull(awsScriptDetails.ExecutablePath),
		ExecutableParameters: util.StringOrNull(awsScriptDetails.ExecutableParameters),
		TimeoutInSeconds:     util.Int32OrNull(awsScriptDetails.TimeoutInSeconds),
	}

	if awsScriptDetails.ScriptS3Location != nil {
		detailsModel.ScriptS3Location = flattenSourceS3LocationData(ctx, awsScriptDetails.ScriptS3Location, diags)
	} else {
		detailsModel.ScriptS3Location = types.ObjectNull(s3LocationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(ctx, scriptDetailsObjectType.AttrTypes, detailsModel)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(scriptDetailsObjectType.AttrTypes)
	}

	return obj
}

func flattenAppBlockErrorsData(ctx context.Context, awsAppBlockErrors []awstypes.ErrorDetails, diags *diag.Diagnostics) types.Set {
	if len(awsAppBlockErrors) == 0 {
		return types.SetNull(errorObjectType)
	}

	out := make([]appBlockErrorModel, 0, len(awsAppBlockErrors))

	for _, e := range awsAppBlockErrors {
		out = append(out, appBlockErrorModel{
			ErrorCode:    util.StringOrNull(e.ErrorCode),
			ErrorMessage: util.StringOrNull(e.ErrorMessage),
		})
	}

	setVal, d := types.SetValueFrom(ctx, errorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(errorObjectType)
	}

	return setVal
}
