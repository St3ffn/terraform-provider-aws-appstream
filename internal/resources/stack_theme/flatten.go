// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var s3LocationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"s3_bucket": types.StringType,
		"s3_key":    types.StringType,
	},
}

var footerLinkObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"display_name":    types.StringType,
		"footer_link_url": types.StringType,
	},
}

func flattenS3Location(
	ctx context.Context, awsS3Location *awstypes.S3Location, diags *diag.Diagnostics,
) types.Object {

	if awsS3Location == nil {
		return types.ObjectNull(s3LocationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		s3LocationObjectType.AttrTypes,
		resourceModelOrganizationLogoS3Location{
			S3Bucket: util.StringOrNull(awsS3Location.S3Bucket),
			S3Key:    util.StringOrNull(awsS3Location.S3Key),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(s3LocationObjectType.AttrTypes)
	}

	return obj
}

func flattenFooterLinks(
	ctx context.Context, awsThemeFooterLinks []awstypes.ThemeFooterLink, diags *diag.Diagnostics,
) types.Set {

	if len(awsThemeFooterLinks) == 0 {
		return types.SetNull(footerLinkObjectType)
	}

	out := make([]resourceModelFooterLinks, 0, len(awsThemeFooterLinks))
	for _, link := range awsThemeFooterLinks {
		out = append(out, resourceModelFooterLinks{
			DisplayName:   util.StringOrNull(link.DisplayName),
			FooterLinkUrl: util.StringOrNull(link.FooterLinkURL),
		})
	}

	setVal, d := types.SetValueFrom(ctx, footerLinkObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(footerLinkObjectType)
	}

	return setVal
}
