// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandIconS3Location(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.S3Location {
	var m iconS3LocationModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.S3Location{
		S3Bucket: util.StringPointerOrNil(m.S3Bucket),
		S3Key:    util.StringPointerOrNil(m.S3Key),
	}
}

func expandPlatforms(ctx context.Context, set types.Set, diags *diag.Diagnostics) []awstypes.PlatformType {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var values []string
	diags.Append(set.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil
	}

	if len(values) == 0 {
		return nil
	}

	out := make([]awstypes.PlatformType, 0, len(values))
	for _, v := range values {
		out = append(out, awstypes.PlatformType(v))
	}

	return out
}
