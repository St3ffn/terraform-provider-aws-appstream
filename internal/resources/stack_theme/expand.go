// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandS3Location(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.S3Location {
	var m resourceModelOrganizationLogoS3Location
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.S3Location{
		S3Bucket: util.StringPointerOrNil(m.S3Bucket),
		S3Key:    util.StringPointerOrNil(m.S3Key),
	}
}

func expandFooterLinks(ctx context.Context, set types.Set, diags *diag.Diagnostics) []awstypes.ThemeFooterLink {
	var models []resourceModelFooterLinks
	diags.Append(set.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	if len(models) == 0 {
		return nil
	}

	out := make([]awstypes.ThemeFooterLink, 0, len(models))
	for _, m := range models {
		out = append(out, awstypes.ThemeFooterLink{
			DisplayName:   util.StringPointerOrNil(m.DisplayName),
			FooterLinkURL: util.StringPointerOrNil(m.FooterLinkURL),
		})
	}

	return out
}
