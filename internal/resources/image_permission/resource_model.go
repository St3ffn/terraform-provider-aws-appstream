// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_permission

import "github.com/hashicorp/terraform-plugin-framework/types"

type resourceModel struct {
	// ID is a synthetic identifier composed of "<name>|<shared_account_id>" (computed).
	ID types.String `tfsdk:"id"`
	// Name is the name of the private AppStream image (required).
	Name types.String `tfsdk:"name"`
	// SharedAccountID is the AWS account ID with which the image is shared (required).
	SharedAccountID types.String `tfsdk:"shared_account_id"`
	// ImagePermissions defines how the shared AWS account is allowed to use the image (required).
	ImagePermissions types.Object `tfsdk:"image_permissions"`
}

type imagePermissionsModel struct {
	// AllowFleet indicates whether the image can be used to create or update AppStream fleets
	// in the shared AWS account (required).
	AllowFleet types.Bool `tfsdk:"allow_fleet"`
	// AllowImageBuilder indicates whether the image can be used to create AppStream image builders
	// in the shared AWS account (required).
	AllowImageBuilder types.Bool `tfsdk:"allow_image_builder"`
}
