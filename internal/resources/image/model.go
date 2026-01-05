// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image

import "github.com/hashicorp/terraform-plugin-framework/types"

type model struct {
	// ID is the identifier of the image.
	// This is equal to the image ARN (computed).
	ID types.String `tfsdk:"id"`
	// ARN is the ARN of the AppStream image.
	// Used as lookup input or populated from AWS (optional, computed).
	ARN types.String `tfsdk:"arn"`
	// Name is the name of the AppStream image.
	// Cannot be used together with ARN or NameRegex (optional, computed).
	Name types.String `tfsdk:"name"`
	// NameRegex is a regular expression used to match image names.
	// Cannot be used together with ARN or Name (optional).
	NameRegex types.String `tfsdk:"name_regex"`
	// Visibility filters images by visibility.
	// Valid values are PUBLIC, PRIVATE, or SHARED (optional, computed).
	Visibility types.String `tfsdk:"visibility"`
	// MostRecent controls behavior when multiple images match.
	// If true, the most recent image is selected (optional).
	MostRecent types.Bool `tfsdk:"most_recent"`
	// BaseImageARN is the ARN of the image from which this image was created (computed).
	BaseImageARN types.String `tfsdk:"base_image_arn"`
	// DisplayName is the name displayed to users for the image (computed).
	DisplayName types.String `tfsdk:"display_name"`
	// State is the current lifecycle state of the image (computed).
	State types.String `tfsdk:"state"`
	// ImageBuilderSupported indicates whether an image builder can be launched
	// from this image (computed).
	ImageBuilderSupported types.Bool `tfsdk:"image_builder_supported"`
	// ImageBuilderName is the name of the image builder used to create the image,
	// if applicable (computed).
	ImageBuilderName types.String `tfsdk:"image_builder_name"`
	// Platform is the operating system platform of the image (computed).
	Platform types.String `tfsdk:"platform"`
	// Description is the image description, if set (computed).
	Description types.String `tfsdk:"description"`
	// StateChangeReason describes why the image last changed state (computed).
	StateChangeReason types.Object `tfsdk:"state_change_reason"`
	// Applications is the set of applications included in the image (computed).
	Applications types.Set `tfsdk:"applications"`
	// CreatedTime is the timestamp when the image was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// PublicBaseImageReleasedDate is the release date of the public base image
	// used to create this image (computed).
	PublicBaseImageReleasedDate types.String `tfsdk:"public_base_image_released_date"`
	// AppstreamAgentVersion is the AppStream agent version used by the image (computed).
	AppstreamAgentVersion types.String `tfsdk:"appstream_agent_version"`
	// ImagePermissions describes permissions granted for the image (computed).
	ImagePermissions types.Object `tfsdk:"image_permissions"`
	// ImageErrors is the list of errors reported by AWS for the image (computed).
	ImageErrors types.Set `tfsdk:"image_errors"`
	// LatestAppstreamAgentVersion indicates whether the image uses the latest
	// AppStream agent version (computed).
	LatestAppstreamAgentVersion types.String `tfsdk:"latest_appstream_agent_version"`
	// SupportedInstanceFamilies lists the instance families supported by the image (computed).
	SupportedInstanceFamilies types.Set `tfsdk:"supported_instance_families"`
	// DynamicAppProvidersEnabled indicates whether dynamic app providers
	// are enabled for the image (computed).
	DynamicAppProvidersEnabled types.String `tfsdk:"dynamic_app_providers_enabled"`
	// ImageSharedWithOthers indicates whether the image is shared with other AWS accounts (computed).
	ImageSharedWithOthers types.String `tfsdk:"image_shared_with_others"`
	// ManagedSoftwareIncluded indicates whether the image includes managed software (computed).
	ManagedSoftwareIncluded types.Bool `tfsdk:"managed_software_included"`
	// ImageType is the type of the image: CUSTOM or NATIVE (computed).
	ImageType types.String `tfsdk:"image_type"`
	// Tags is a map of tags assigned to the image (computed).
	Tags types.Map `tfsdk:"tags"`
}

type stateChangeReasonModel struct {
	// Code is the state change reason code (computed).
	Code types.String `tfsdk:"code"`
	// Message is the human-readable state change reason message (computed).
	Message types.String `tfsdk:"message"`
}

type applicationModel struct {
	// Name is the application name (computed).
	Name types.String `tfsdk:"name"`
	// DisplayName is the application display name (computed).
	DisplayName types.String `tfsdk:"display_name"`
	// IconURL is the URL of the application icon (computed).
	IconURL types.String `tfsdk:"icon_url"`
	// LaunchPath is the path to the application executable (computed).
	LaunchPath types.String `tfsdk:"launch_path"`
	// LaunchParameters are the parameters passed at application launch (computed).
	LaunchParameters types.String `tfsdk:"launch_parameters"`
	// Enabled indicates whether the application is enabled (computed).
	Enabled types.Bool `tfsdk:"enabled"`
	// Metadata is additional application metadata (computed).
	Metadata types.Map `tfsdk:"metadata"`
	// WorkingDirectory is the application working directory (computed).
	WorkingDirectory types.String `tfsdk:"working_directory"`
	// Description is the application description (computed).
	Description types.String `tfsdk:"description"`
	// ARN is the ARN of the application (computed).
	ARN types.String `tfsdk:"arn"`
	// AppBlockARN is the ARN of the associated app block (computed).
	AppBlockARN types.String `tfsdk:"app_block_arn"`
	// IconS3Location is the S3 location of the application icon (computed).
	IconS3Location types.Object `tfsdk:"icon_s3_location"`
	// Platforms are the platforms on which the application can run (computed).
	Platforms types.Set `tfsdk:"platforms"`
	// InstanceFamilies are the instance families supported by the application (computed).
	InstanceFamilies types.Set `tfsdk:"instance_families"`
	// CreatedTime is the timestamp when the application was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
}

type iconS3LocationModel struct {
	// S3Bucket is the S3 bucket name (computed).
	S3Bucket types.String `tfsdk:"s3_bucket"`
	// S3Key is the S3 object key (computed).
	S3Key types.String `tfsdk:"s3_key"`
}

type imagePermissionsModel struct {
	// AllowFleet indicates whether the image can be used by fleets (computed).
	AllowFleet types.Bool `tfsdk:"allow_fleet"`
	// AllowImageBuilder indicates whether the image can be used by image builders (computed).
	AllowImageBuilder types.Bool `tfsdk:"allow_image_builder"`
}

type imageErrorModel struct {
	// ErrorCode is the error code reported by AWS (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
	// ErrorTimestamp is the timestamp when the error occurred (computed).
	ErrorTimestamp types.String `tfsdk:"error_timestamp"`
}
