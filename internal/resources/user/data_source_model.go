// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// ID is a synthetic identifier composed of "<authentication_type>|<user_name>" (computed).
	ID types.String `tfsdk:"id"`
	// AuthenticationType is the authentication type for the user (required).
	AuthenticationType types.String `tfsdk:"authentication_type"`
	// UserName is the email address of the AppStream user (required).
	UserName types.String `tfsdk:"user_name"`
	// FirstName is the first (given) name of the user (computed).
	FirstName types.String `tfsdk:"first_name"`
	// LastName is the last (family) name of the user (computed).
	LastName types.String `tfsdk:"last_name"`
	// Enabled specifies whether the user is enabled (computed).
	Enabled types.Bool `tfsdk:"enabled"`
	// Status is the current status of the user as reported by AWS (computed).
	Status types.String `tfsdk:"status"`
	// ARN is the Amazon Resource Name of the AppStream user (computed).
	ARN types.String `tfsdk:"arn"`
	// CreatedTime is the timestamp when the user was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
}
