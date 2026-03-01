// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package export_image_tasks

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandFilters(ctx context.Context, set types.Set, diags *diag.Diagnostics) []awstypes.Filter {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var models []dataSourceModelFilters
	diags.Append(set.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	if len(models) == 0 {
		return nil
	}

	out := make([]awstypes.Filter, 0, len(models))
	for _, m := range models {
		out = append(out, awstypes.Filter{
			Name:   util.StringPointerOrNil(m.Name),
			Values: util.ExpandStringSetOrNil(ctx, m.Values, diags),
		})
	}

	return out
}
