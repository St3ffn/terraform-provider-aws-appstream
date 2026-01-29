// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package application

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// ID is the ARN of the AppStream application.
	ID types.String `tfsdk:"id"`
	// ARN is the ARN of the AppStream application (required).
	ARN types.String `tfsdk:"arn"`
	// Name is the name of the AppStream application (computed).
	Name types.String `tfsdk:"name"`
	// DisplayName is the name of the application as displayed to users (computed).
	DisplayName types.String `tfsdk:"display_name"`
	// Description is the description of the application (computed).
	Description types.String `tfsdk:"description"`
	// IconS3Location specifies the S3 location of the application icon (computed).
	IconS3Location types.Object `tfsdk:"icon_s3_location"`
	// LaunchPath is the path to the application executable within the image (computed).
	LaunchPath types.String `tfsdk:"launch_path"`
	// WorkingDirectory is the working directory of the application (computed).
	WorkingDirectory types.String `tfsdk:"working_directory"`
	// LaunchParameters are the parameters passed to the application at launch (computed).
	LaunchParameters types.String `tfsdk:"launch_parameters"`
	// Platforms specifies the platforms the application supports (computed).
	Platforms types.Set `tfsdk:"platforms"`
	// InstanceFamilies specifies the instance families the application supports (computed).
	InstanceFamilies types.Set `tfsdk:"instance_families"`
	// AppBlockARN is the ARN of the app block associated with the application (computed).
	AppBlockARN types.String `tfsdk:"app_block_arn"`
	// Tags is a map of tags to assign to the application (computed).
	Tags types.Map `tfsdk:"tags"`
	// CreatedTime is the timestamp when the application was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
}
