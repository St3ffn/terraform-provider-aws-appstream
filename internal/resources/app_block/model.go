// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import "github.com/hashicorp/terraform-plugin-framework/types"

type model struct {
	// ID is the ARN of the AppStream app block.
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream app block (required).
	Name types.String `tfsdk:"name"`
	// DisplayName is the display name of the app block (optional).
	DisplayName types.String `tfsdk:"display_name"`
	// Description is the description of the app block (optional).
	Description types.String `tfsdk:"description"`
	// SourceS3Location specifies the source S3 location of the app block (required).
	SourceS3Location types.Object `tfsdk:"source_s3_location"`
	// SetupScriptDetails specifies the setup script configuration (optional).
	SetupScriptDetails types.Object `tfsdk:"setup_script_details"`
	// PostSetupScriptDetails specifies the post-setup script configuration (optional).
	PostSetupScriptDetails types.Object `tfsdk:"post_setup_script_details"`
	// PackagingType specifies the packaging type of the app block (optional).
	PackagingType types.String `tfsdk:"packaging_type"`
	// Tags is a map of tags assigned to the app block (optional).
	Tags types.Map `tfsdk:"tags"`
	// ARN is the ARN of the AppStream app block (computed).
	ARN types.String `tfsdk:"arn"`
	// CreatedTime is the timestamp when the app block was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// AppBlockErrors is the list of errors reported by AWS for the app block (computed).
	AppBlockErrors types.Set `tfsdk:"app_block_errors"`
}

type sourceS3LocationModel struct {
	// S3Bucket is the name of the Amazon S3 bucket (required).
	S3Bucket types.String `tfsdk:"s3_bucket"`
	// S3Key is the S3 object key (optional).
	S3Key types.String `tfsdk:"s3_key"`
}

type scriptDetailsModel struct {
	// ScriptS3Location specifies the S3 location of the script (required).
	ScriptS3Location types.Object `tfsdk:"script_s3_location"`
	// ExecutablePath is the run path for the script (required).
	ExecutablePath types.String `tfsdk:"executable_path"`
	// ExecutableParameters are the runtime parameters passed to the script (optional).
	ExecutableParameters types.String `tfsdk:"executable_parameters"`
	// TimeoutInSeconds is the timeout for the script execution (required).
	TimeoutInSeconds types.Int32 `tfsdk:"timeout_in_seconds"`
}

type appBlockErrorModel struct {
	// ErrorCode is the error code reported by AWS (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
}
