// Copyright (c) St3ffn
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

var imagePermissionEntryObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"shared_account_id": types.StringType,
		"image_permissions": imagePermissionsObjectType,
	},
}

func flattenImagePermissionsData(
	ctx context.Context, awsEntries []awstypes.SharedImagePermissions, diags *diag.Diagnostics,
) types.Set {

	if len(awsEntries) == 0 {
		return types.SetNull(imagePermissionEntryObjectType)
	}

	entries := make([]imagePermissionModel, 0, len(awsEntries))

	for _, e := range awsEntries {
		entries = append(entries, imagePermissionModel{
			SharedAccountID:  util.StringOrNull(e.SharedAccountId),
			ImagePermissions: flattenImagePermissionsResource(ctx, e.ImagePermissions, diags),
		})
	}

	setVal, d := types.SetValueFrom(ctx, imagePermissionEntryObjectType, entries)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(imagePermissionEntryObjectType)
	}

	return setVal
}
