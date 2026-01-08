// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// ID is a synthetic identifier composed of "<name>".
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream image builder (required).
	Name types.String `tfsdk:"name"`
	// ImageARN is the ARN of the AppStream image used to create the image builder (computed).
	ImageARN types.String `tfsdk:"image_arn"`
	// InstanceType is the instance type used to launch the image builder (computed).
	InstanceType types.String `tfsdk:"instance_type"`
	// Description is a description to display for the image builder (computed).
	Description types.String `tfsdk:"description"`
	// DisplayName is the name of the image builder shown to users (computed).
	DisplayName types.String `tfsdk:"display_name"`
	// VPCConfig specifies the VPC configuration for the image builder (computed).
	VPCConfig types.Object `tfsdk:"vpc_config"`
	// IAMRoleARN is the ARN of the IAM role applied to the image builder (computed).
	IAMRoleARN types.String `tfsdk:"iam_role_arn"`
	// EnableDefaultInternetAccess specifies whether the image builder has internet access (computed).
	EnableDefaultInternetAccess types.Bool `tfsdk:"enable_default_internet_access"`
	// DomainJoinInfo specifies Active Directory domain join configuration (computed).
	DomainJoinInfo types.Object `tfsdk:"domain_join_info"`
	// AppstreamAgentVersion is the AppStream agent version used by the image builder (computed).
	AppstreamAgentVersion types.String `tfsdk:"appstream_agent_version"`
	// AccessEndpoints specifies interface VPC endpoints used to access the image builder (computed).
	AccessEndpoints types.Set `tfsdk:"access_endpoints"`
	// RootVolumeConfig specifies the root volume configuration for the image builder (computed).
	RootVolumeConfig types.Object `tfsdk:"root_volume_config"`
	// Tags is a map of tags assigned to the image builder (computed).
	Tags types.Map `tfsdk:"tags"`
	// ARN is the ARN of the AppStream image builder (computed).
	ARN types.String `tfsdk:"arn"`
	// CreatedTime is the timestamp when the image builder was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// Platform is the operating system platform of the image builder (computed).
	Platform types.String `tfsdk:"platform"`
	// NetworkAccessConfiguration contains network details of the image builder (computed).
	NetworkAccessConfiguration types.Object `tfsdk:"network_access_configuration"`
	// LatestAppstreamAgentVersion indicates whether the latest AppStream agent is used (computed).
	LatestAppstreamAgentVersion types.String `tfsdk:"latest_appstream_agent_version"`
	// State is the current state of the image builder (computed).
	State types.String `tfsdk:"state"`
	// StateChangeReason describes the most recent state change, if any (computed).
	StateChangeReason types.Object `tfsdk:"state_change_reason"`
	// ImageBuilderErrors is the list of errors reported by AWS for the image builder (computed).
	ImageBuilderErrors types.Set `tfsdk:"image_builder_errors"`
}
