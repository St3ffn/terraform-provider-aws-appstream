// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// ID is a synthetic identifier composed of "<name>".
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream fleet (required).
	Name types.String `tfsdk:"name"`
	// ImageName is the name of the image used to create the fleet.
	// Exactly one of ImageName or ImageARN must be specified (computed).
	ImageName types.String `tfsdk:"image_name"`
	// ImageARN is the ARN of the image used to create the fleet.
	// Exactly one of ImageName or ImageARN must be specified (computed).
	ImageARN types.String `tfsdk:"image_arn"`
	// InstanceType is the EC2 instance type for the fleet (computed).
	InstanceType types.String `tfsdk:"instance_type"`
	// FleetType is the type of fleet: ON_DEMAND, ALWAYS_ON, or ELASTIC (computed).
	FleetType types.String `tfsdk:"fleet_type"`
	// ComputeCapacity specifies the desired number of instances or sessions (computed).
	ComputeCapacity types.Object `tfsdk:"compute_capacity"`
	// VPCConfig specifies the VPC configuration for the fleet (computed).
	VPCConfig types.Object `tfsdk:"vpc_config"`
	// MaxUserDurationInSeconds is the maximum streaming session length (computed).
	MaxUserDurationInSeconds types.Int32 `tfsdk:"max_user_duration_in_seconds"`
	// DisconnectTimeoutInSeconds is the time before a disconnected session is terminated (computed).
	DisconnectTimeoutInSeconds types.Int32 `tfsdk:"disconnect_timeout_in_seconds"`
	// IdleDisconnectTimeoutInSeconds is the timeout for idle streaming sessions (computed).
	IdleDisconnectTimeoutInSeconds types.Int32 `tfsdk:"idle_disconnect_timeout_in_seconds"`
	// Description is a description to display for the fleet (computed).
	Description types.String `tfsdk:"description"`
	// DisplayName is the fleet name shown to users (computed).
	DisplayName types.String `tfsdk:"display_name"`
	// DisableIMDSV1 specifies whether the fleet has disabled IMDSv1 and IMDSv2 is enforced (computed).
	DisableIMDSV1 types.Bool `tfsdk:"disable_imds_v1"`
	// EnableDefaultInternetAccess enables outbound internet access (computed).
	EnableDefaultInternetAccess types.Bool `tfsdk:"enable_default_internet_access"`
	// DomainJoinInfo specifies Active Directory domain join configuration (computed).
	DomainJoinInfo types.Object `tfsdk:"domain_join_info"`
	// IAMRoleARN is the ARN of the IAM role applied to the fleet instances (computed).
	IAMRoleARN types.String `tfsdk:"iam_role_arn"`
	// StreamView controls which streaming protocol views are enabled (computed).
	StreamView types.String `tfsdk:"stream_view"`
	// Platform is the platform type of the fleet (computed).
	Platform types.String `tfsdk:"platform"`
	// MaxConcurrentSessions is the maximum number of concurrent streaming sessions (computed).
	MaxConcurrentSessions types.Int32 `tfsdk:"max_concurrent_sessions"`
	// MaxSessionsPerInstance is the maximum number of user sessions allowed per fleet instance.
	// This setting applies only to multi-session fleets (computed).
	MaxSessionsPerInstance types.Int32 `tfsdk:"max_sessions_per_instance"`
	// USBDeviceFilterStrings defines which USB devices are allowed (computed).
	USBDeviceFilterStrings types.Set `tfsdk:"usb_device_filter_strings"`
	// SessionScriptS3Location specifies the S3 location of the session scripts
	// configuration ZIP file. This setting applies only to elastic fleets. (computed).
	SessionScriptS3Location types.Object `tfsdk:"session_script_s3_location"`
	// RootVolumeConfig specifies the root volume configuration for the fleet (computed).
	RootVolumeConfig types.Object `tfsdk:"root_volume_config"`
	// Tags is a map of tags to assign to the fleet (computed).
	Tags types.Map `tfsdk:"tags"`
	// ARN is the ARN of the AppStream fleet (computed).
	ARN types.String `tfsdk:"arn"`
	// CreatedTime is the timestamp when the fleet was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// State is the state of the AppStream fleet (computed).
	State types.String `tfsdk:"state"`
	// FleetErrors is the list of errors reported by AWS for the fleet (computed).
	FleetErrors types.Set `tfsdk:"fleet_errors"`
}
