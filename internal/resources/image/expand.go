// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandStates(ctx context.Context, set types.Set, diags *diag.Diagnostics) map[awstypes.ImageState]struct{} {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	values := util.ExpandStringSetOrNil(ctx, set, diags)
	if diags.HasError() {
		return nil
	}

	if len(values) == 0 {
		return nil
	}

	out := make(map[awstypes.ImageState]struct{}, len(values))
	for _, v := range values {
		out[awstypes.ImageState(v)] = struct{}{}
	}

	return out
}
