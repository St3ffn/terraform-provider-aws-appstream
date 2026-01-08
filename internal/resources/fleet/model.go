// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

import "github.com/hashicorp/terraform-plugin-framework/types"

type model struct {
	// ID is a synthetic identifier composed of "<name>".
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream fleet (required).
	Name types.String `tfsdk:"name"`
	// ImageName is the name of the image used to create the fleet.
	// Exactly one of ImageName or ImageARN must be specified (optional).
	ImageName types.String `tfsdk:"image_name"`
	// ImageARN is the ARN of the image used to create the fleet.
	// Exactly one of ImageName or ImageARN must be specified (optional).
	ImageARN types.String `tfsdk:"image_arn"`
	// InstanceType is the EC2 instance type for the fleet (required).
	InstanceType types.String `tfsdk:"instance_type"`
	// FleetType is the type of fleet: ON_DEMAND, ALWAYS_ON, or ELASTIC (required).
	FleetType types.String `tfsdk:"fleet_type"`
	// ComputeCapacity specifies the desired number of instances or sessions (non-elastic fleets only).
	ComputeCapacity types.Object `tfsdk:"compute_capacity"`
	// VPCConfig specifies the VPC configuration for the fleet (required for elastic fleets).
	VPCConfig types.Object `tfsdk:"vpc_config"`
	// MaxUserDurationInSeconds is the maximum streaming session length (optional, computed).
	MaxUserDurationInSeconds types.Int32 `tfsdk:"max_user_duration_in_seconds"`
	// DisconnectTimeoutInSeconds is the time before a disconnected session is terminated (optional, computed).
	DisconnectTimeoutInSeconds types.Int32 `tfsdk:"disconnect_timeout_in_seconds"`
	// IdleDisconnectTimeoutInSeconds is the timeout for idle streaming sessions (optional, computed).
	IdleDisconnectTimeoutInSeconds types.Int32 `tfsdk:"idle_disconnect_timeout_in_seconds"`
	// Description is a description to display for the fleet (optional).
	Description types.String `tfsdk:"description"`
	// DisplayName is the fleet name shown to users (optional).
	DisplayName types.String `tfsdk:"display_name"`
	// EnableDefaultInternetAccess enables outbound internet access (optional).
	EnableDefaultInternetAccess types.Bool `tfsdk:"enable_default_internet_access"`
	// DomainJoinInfo specifies Active Directory domain join configuration (Windows fleets only).
	DomainJoinInfo types.Object `tfsdk:"domain_join_info"`
	// IAMRoleARN is the ARN of the IAM role applied to the fleet instances (optional).
	IAMRoleARN types.String `tfsdk:"iam_role_arn"`
	// StreamView controls which streaming protocol views are enabled (optional).
	StreamView types.String `tfsdk:"stream_view"`
	// Platform is the platform type of the fleet (optional).
	Platform types.String `tfsdk:"platform"`
	// MaxConcurrentSessions is the maximum number of concurrent streaming sessions. (elastic fleets only).
	MaxConcurrentSessions types.Int32 `tfsdk:"max_concurrent_sessions"`
	// MaxSessionsPerInstance is the maximum number of user sessions allowed per fleet instance.
	// This setting applies only to multi-session fleets (optional).
	MaxSessionsPerInstance types.Int32 `tfsdk:"max_sessions_per_instance"`
	// USBDeviceFilterStrings defines which USB devices are allowed (Windows fleets only).
	USBDeviceFilterStrings types.Set `tfsdk:"usb_device_filter_strings"`
	// SessionScriptS3Location specifies the S3 location of the session scripts
	// configuration ZIP file. This setting applies only to elastic fleets. (optional).
	SessionScriptS3Location types.Object `tfsdk:"session_script_s3_location"`
	// RootVolumeConfig specifies the root volume configuration for the fleet (optional).
	RootVolumeConfig types.Object `tfsdk:"root_volume_config"`
	// Tags is a map of tags to assign to the fleet (optional).
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

type computeCapacityModel struct {
	// DesiredInstances is the desired number of streaming instances for a
	// non-elastic fleet. This must be specified for single-session fleets and
	// cannot be used together with DesiredSessions.
	DesiredInstances types.Int32 `tfsdk:"desired_instances"`
	// DesiredSessions is the desired number of concurrent user sessions for a
	// non-elastic multi-session fleet. This must be specified for multi-session
	// fleets and cannot be used together with DesiredInstances.
	DesiredSessions types.Int32 `tfsdk:"desired_sessions"`
}

type vpcConfigModel struct {
	// SubnetIDs are the subnet IDs for the fleet.
	// Required for elastic fleets. At least two subnets in different Availability Zones
	// must be specified. The Availability Zone requirement is enforced by AWS.
	SubnetIDs types.Set `tfsdk:"subnet_ids"`
	// SecurityGroupIDs are the security group IDs for the fleet.
	SecurityGroupIDs types.Set `tfsdk:"security_group_ids"`
}

type domainJoinInfoModel struct {
	// DirectoryName is the name of the Active Directory.
	DirectoryName types.String `tfsdk:"directory_name"`
	// OrganizationalUnitDistinguishedName is the OU DN for computer accounts.
	OrganizationalUnitDistinguishedName types.String `tfsdk:"organizational_unit_distinguished_name"`
}

type sessionScriptS3LocationModel struct {
	// S3Bucket is the S3 bucket containing the session script.
	S3Bucket types.String `tfsdk:"s3_bucket"`
	// S3Key is the S3 object key of the session script.
	S3Key types.String `tfsdk:"s3_key"`
}

type rootVolumeConfigModel struct {
	// VolumeSizeInGB is the size of the root volume.
	VolumeSizeInGB types.Int32 `tfsdk:"volume_size_in_gb"`
}

type fleetErrorModel struct {
	// ErrorCode is the error code reported by AWS (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
}
