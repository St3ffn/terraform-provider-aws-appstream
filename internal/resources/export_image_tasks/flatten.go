// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package export_image_tasks

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var exportImageTaskObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"task_id":            types.StringType,
		"image_arn":          types.StringType,
		"ami_name":           types.StringType,
		"ami_description":    types.StringType,
		"ami_id":             types.StringType,
		"state":              types.StringType,
		"created_date":       types.StringType,
		"tag_specifications": types.MapType{ElemType: types.StringType},
		"error_details":      types.SetType{ElemType: errorDetailObjectType},
	},
}

var errorDetailObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenExportImageTasksData(
	ctx context.Context, tasks []awstypes.ExportImageTask, diags *diag.Diagnostics,
) types.Set {

	if len(tasks) == 0 {
		return types.SetNull(exportImageTaskObjectType)
	}

	out := make([]exportImageTaskModel, 0, len(tasks))

	for _, t := range tasks {
		out = append(out, exportImageTaskModel{
			TaskID:            util.StringOrNull(t.TaskId),
			ImageArn:          util.StringOrNull(t.ImageArn),
			AmiName:           util.StringOrNull(t.AmiName),
			AmiDescription:    util.StringOrNull(t.AmiDescription),
			AmiID:             util.StringOrNull(t.AmiId),
			State:             types.StringValue(string(t.State)),
			CreatedDate:       util.StringFromTime(t.CreatedDate),
			TagSpecifications: util.MapStringOrNull(ctx, t.TagSpecifications, diags),
			ErrorDetails:      flattenErrorDetailsData(ctx, t.ErrorDetails, diags),
		})
	}

	setVal, d := types.SetValueFrom(ctx, exportImageTaskObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(exportImageTaskObjectType)
	}

	return setVal
}

func flattenErrorDetailsData(
	ctx context.Context, errorDetails []awstypes.ErrorDetails, diags *diag.Diagnostics,
) types.Set {

	if len(errorDetails) == 0 {
		return types.SetNull(errorDetailObjectType)
	}

	out := make([]errorDetailModel, 0, len(errorDetails))

	for _, e := range errorDetails {
		out = append(out, errorDetailModel{
			ErrorCode:    util.StringOrNull(e.ErrorCode),
			ErrorMessage: util.StringOrNull(e.ErrorMessage),
		})
	}

	setVal, d := types.SetValueFrom(ctx, errorDetailObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(errorDetailObjectType)
	}

	return setVal
}
