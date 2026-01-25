// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// ID is the ARN of the AppStream app block.
	ID types.String `tfsdk:"id"`
	// ARN is the ARN of the AppStream app block (required).
	ARN types.String `tfsdk:"arn"`
	// Name is the name of the AppStream app block (computed).
	Name types.String `tfsdk:"name"`
	// DisplayName is the display name of the app block (computed).
	DisplayName types.String `tfsdk:"display_name"`
	// Description is the description of the app block (computed).
	Description types.String `tfsdk:"description"`
	// SourceS3Location specifies the source S3 location of the app block (computed).
	SourceS3Location types.Object `tfsdk:"source_s3_location"`
	// SetupScriptDetails specifies the setup script configuration (computed).
	SetupScriptDetails types.Object `tfsdk:"setup_script_details"`
	// PostSetupScriptDetails specifies the post-setup script configuration (computed).
	PostSetupScriptDetails types.Object `tfsdk:"post_setup_script_details"`
	// PackagingType specifies the packaging type of the app block (computed).
	PackagingType types.String `tfsdk:"packaging_type"`
	// Tags is a map of tags assigned to the app block (computed).
	Tags types.Map `tfsdk:"tags"`
	// CreatedTime is the timestamp when the app block was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// State is the state of the AppStream app block (computed).
	State types.String `tfsdk:"state"`
	// AppBlockErrors is the list of errors reported by AWS for the app block (computed).
	AppBlockErrors types.Set `tfsdk:"app_block_errors"`
}
