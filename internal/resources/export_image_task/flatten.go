// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package export_image_task

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var errorDetailObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenErrorDetailsData(ctx context.Context, errorDetails []awstypes.ErrorDetails, diags *diag.Diagnostics) types.Set {
	if len(errorDetails) == 0 {
		return types.SetNull(errorDetailObjectType)
	}

	out := make([]dataSourceModelErrorDetails, 0, len(errorDetails))

	for _, e := range errorDetails {
		out = append(out, dataSourceModelErrorDetails{
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
