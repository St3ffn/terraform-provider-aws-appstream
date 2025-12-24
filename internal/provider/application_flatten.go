// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var applicationIconS3LocationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"s3_bucket": types.StringType,
		"s3_key":    types.StringType,
	},
}

func flattenApplicationIconS3Location(
	ctx context.Context,
	awsS3Location *awstypes.S3Location,
	diags *diag.Diagnostics,
) types.Object {

	if awsS3Location == nil {
		return types.ObjectNull(applicationIconS3LocationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		applicationIconS3LocationObjectType.AttrTypes,
		applicationIconS3LocationModel{
			S3Bucket: stringOrNull(awsS3Location.S3Bucket),
			S3Key:    stringOrNull(awsS3Location.S3Key),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(applicationIconS3LocationObjectType.AttrTypes)
	}

	return obj
}

func flattenApplicationPlatforms(
	ctx context.Context,
	platforms []awstypes.PlatformType,
	diags *diag.Diagnostics,
) types.Set {

	if len(platforms) == 0 {
		return types.SetValueMust(types.StringType, []attr.Value{})
	}

	values := make([]string, 0, len(platforms))
	for _, p := range platforms {
		values = append(values, string(p))
	}

	setVal, d := types.SetValueFrom(ctx, types.StringType, values)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(types.StringType)
	}

	return setVal
}
