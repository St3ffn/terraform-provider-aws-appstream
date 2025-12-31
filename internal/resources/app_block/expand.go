// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandSourceS3Location(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.S3Location {
	var m sourceS3LocationModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.S3Location{
		S3Bucket: util.StringPointerOrNil(m.S3Bucket),
		S3Key:    util.StringPointerOrNil(m.S3Key),
	}
}

func expandScriptDetails(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.ScriptDetails {
	var m scriptDetailsModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	script := &awstypes.ScriptDetails{
		ExecutablePath:       util.StringPointerOrNil(m.ExecutablePath),
		ExecutableParameters: util.StringPointerOrNil(m.ExecutableParameters),
		TimeoutInSeconds:     util.Int32PointerOrNil(m.TimeoutInSeconds),
	}

	if !m.ScriptS3Location.IsNull() && !m.ScriptS3Location.IsUnknown() {
		script.ScriptS3Location = expandSourceS3Location(ctx, m.ScriptS3Location, diags)
	}

	return script
}
