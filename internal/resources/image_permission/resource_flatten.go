// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_permission

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var imagePermissionsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"allow_fleet":         types.BoolType,
		"allow_image_builder": types.BoolType,
	},
}

func flattenImagePermissionsResource(
	ctx context.Context, awsImagePermissions *awstypes.ImagePermissions, diags *diag.Diagnostics,
) types.Object {

	if awsImagePermissions == nil {
		return types.ObjectNull(imagePermissionsObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		imagePermissionsObjectType.AttrTypes,
		resourceModelImagePermissions{
			AllowFleet:        util.BoolOrNull(awsImagePermissions.AllowFleet),
			AllowImageBuilder: util.BoolOrNull(awsImagePermissions.AllowImageBuilder),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(imagePermissionsObjectType.AttrTypes)
	}

	return obj
}
