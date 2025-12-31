// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var s3LocationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"s3_bucket": types.StringType,
		"s3_key":    types.StringType,
	},
}

var scriptDetailsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"script_s3_location":    s3LocationObjectType,
		"executable_path":       types.StringType,
		"executable_parameters": types.StringType,
		"timeout_in_seconds":    types.Int32Type,
	},
}

var sourceS3LocationObjectType = s3LocationObjectType

var errorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenScriptDetailsResource(
	ctx context.Context, prior types.Object, awsScriptDetails *awstypes.ScriptDetails, diags *diag.Diagnostics,
) types.Object {

	// user never managed this attribute
	if prior.IsNull() {
		return types.ObjectNull(scriptDetailsObjectType.AttrTypes)
	}

	// terraform does not yet know during planning
	if prior.IsUnknown() {
		return types.ObjectUnknown(scriptDetailsObjectType.AttrTypes)
	}

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
