// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_permission

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandImagePermissions(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.ImagePermissions {
	var m imagePermissionsModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	imagePermissions := &awstypes.ImagePermissions{
		AllowFleet:        util.BoolPointerOrNil(m.AllowFleet),
		AllowImageBuilder: util.BoolPointerOrNil(m.AllowImageBuilder),
	}

	return imagePermissions
}
