// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var iconS3LocationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"s3_bucket": types.StringType,
		"s3_key":    types.StringType,
	},
}

func flattenIconS3Location(ctx context.Context, awsS3Location *awstypes.S3Location, diags *diag.Diagnostics) types.Object {
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

func flattenPlatforms(ctx context.Context, platforms []awstypes.PlatformType, diags *diag.Diagnostics) types.Set {
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
