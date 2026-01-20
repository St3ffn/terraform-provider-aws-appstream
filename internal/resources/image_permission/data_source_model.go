// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_permission

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// Name is the name of the private AppStream image (required).
	Name types.String `tfsdk:"name"`
	// Permissions is the list of image permissions for shared AWS accounts (computed).
	Permissions types.Set `tfsdk:"permissions"`
}

type imagePermissionEntryModel struct {
	// SharedAccountID is the AWS account ID the image is shared with (computed).
	SharedAccountID types.String `tfsdk:"shared_account_id"`
	// ImagePermissions defines how the shared AWS account is allowed to use the image (computed).
	ImagePermissions types.Object `tfsdk:"image_permissions"`
}
