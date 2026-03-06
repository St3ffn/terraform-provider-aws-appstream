// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandRuntimeValidationConfig(
	ctx context.Context, obj types.Object, diags *diag.Diagnostics,
) *awstypes.RuntimeValidationConfig {

	var m resourceModelRuntimeValidationConfig
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.RuntimeValidationConfig{
		IntendedInstanceType: util.StringPointerOrNil(m.IntendedInstanceType),
	}
}

func expandAppCatalogConfig(
	ctx context.Context, setVal types.Set, diags *diag.Diagnostics,
) []awstypes.ApplicationConfig {

	var models []resourceModelAppCatalogConfig
	diags.Append(setVal.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	if len(models) == 0 {
		return nil
	}

	out := make([]awstypes.ApplicationConfig, 0, len(models))
	for _, m := range models {
		out = append(out, awstypes.ApplicationConfig{
			Name:                 util.StringPointerOrNil(m.Name),
			AbsoluteAppPath:      util.StringPointerOrNil(m.AbsoluteAppPath),
			DisplayName:          util.StringPointerOrNil(m.DisplayName),
			LaunchParameters:     util.StringPointerOrNil(m.LaunchParameters),
			WorkingDirectory:     util.StringPointerOrNil(m.WorkingDirectory),
			AbsoluteIconPath:     util.StringPointerOrNil(m.AbsoluteIconPath),
			AbsoluteManifestPath: util.StringPointerOrNil(m.AbsoluteManifestPath),
		})
	}

	return out
}
