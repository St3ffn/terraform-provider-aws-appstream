// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type model struct {
	// ID is a synthetic identifier composed of "<app_block_builder_name>|<app_block_arn>".
	ID types.String `tfsdk:"id"`
	// AppBlockBuilderName is the name of the AppStream app block builder
	// to associate with the app block (required).
	AppBlockBuilderName types.String `tfsdk:"app_block_builder_name"`
	// AppBlockARN is the ARN of the AppStream app block to associate
	// with the app block builder (required).
	AppBlockARN types.String `tfsdk:"app_block_arn"`
}
