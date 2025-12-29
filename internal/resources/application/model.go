// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

import "github.com/hashicorp/terraform-plugin-framework/types"

type model struct {
	// ID is the ARN of the AppStream application.
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream application (required).
	Name types.String `tfsdk:"name"`
	// DisplayName is the name of the application as displayed to users (optional).
	DisplayName types.String `tfsdk:"display_name"`
	// Description is the description of the application (optional).
	Description types.String `tfsdk:"description"`
	// IconS3Location specifies the S3 location of the application icon (required).
	IconS3Location types.Object `tfsdk:"icon_s3_location"`
	// LaunchPath is the path to the application executable within the image (required).
	LaunchPath types.String `tfsdk:"launch_path"`
	// WorkingDirectory is the working directory of the application (optional).
	WorkingDirectory types.String `tfsdk:"working_directory"`
	// LaunchParameters are the parameters passed to the application at launch (optional).
	LaunchParameters types.String `tfsdk:"launch_parameters"`
	// Platforms specifies the platforms the application supports (required).
	Platforms types.Set `tfsdk:"platforms"`
	// InstanceFamilies specifies the instance families the application supports (required).
	InstanceFamilies types.Set `tfsdk:"instance_families"`
	// AppBlockARN is the ARN of the app block associated with the application (required).
	AppBlockARN types.String `tfsdk:"app_block_arn"`
	// Tags is a map of tags to assign to the application (optional).
	Tags types.Map `tfsdk:"tags"`
	// ARN is the ARN of the AppStream application (computed).
	ARN types.String `tfsdk:"arn"`
	// CreatedTime is the timestamp when the application was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
}

type iconS3LocationModel struct {
	// S3Bucket is the S3 bucket containing the application icon (required).
	S3Bucket types.String `tfsdk:"s3_bucket"`
	// S3Key is the S3 object key of the application icon (required).
	S3Key types.String `tfsdk:"s3_key"`
}
