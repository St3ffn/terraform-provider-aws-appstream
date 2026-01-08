// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import "github.com/hashicorp/terraform-plugin-framework/types"

type resourceModel struct {
	// ID is a synthetic identifier composed of "<name>".
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream image builder (required).
	Name types.String `tfsdk:"name"`
	// ImageName is the name of the AppStream image used to create the image builder.
	// Exactly one of ImageName or ImageARN must be specified (optional).
	ImageName types.String `tfsdk:"image_name"`
	// ImageARN is the ARN of the AppStream image used to create the image builder.
	// Exactly one of ImageName or ImageARN must be specified (optional, computed).
	ImageARN types.String `tfsdk:"image_arn"`
	// InstanceType is the instance type used to launch the image builder (required).
	InstanceType types.String `tfsdk:"instance_type"`
	// Description is a description to display for the image builder (optional).
	Description types.String `tfsdk:"description"`
	// DisplayName is the name of the image builder shown to users (optional).
	DisplayName types.String `tfsdk:"display_name"`
	// VPCConfig specifies the VPC configuration for the image builder (optional).
	VPCConfig types.Object `tfsdk:"vpc_config"`
	// IAMRoleARN is the ARN of the IAM role applied to the image builder (optional).
	IAMRoleARN types.String `tfsdk:"iam_role_arn"`
	// EnableDefaultInternetAccess specifies whether the image builder has internet access (optional, computed).
	EnableDefaultInternetAccess types.Bool `tfsdk:"enable_default_internet_access"`
	// DomainJoinInfo specifies Active Directory domain join configuration (optional).
	DomainJoinInfo types.Object `tfsdk:"domain_join_info"`
	// AppstreamAgentVersion is the AppStream agent version used by the image builder (optional, computed).
	AppstreamAgentVersion types.String `tfsdk:"appstream_agent_version"`
	// AccessEndpoints specifies interface VPC endpoints used to access the image builder (optional).
	AccessEndpoints types.Set `tfsdk:"access_endpoints"`
	// RootVolumeConfig specifies the root volume configuration for the image builder (optional, computed).
	RootVolumeConfig types.Object `tfsdk:"root_volume_config"`
	// Tags is a map of tags assigned to the image builder (optional).
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

type vpcConfigModel struct {
	// SubnetIDs are the subnet IDs in which the image builder is launched.
	// Image builders use one subnet. The constraint is enforced by AWS.
	SubnetIDs types.Set `tfsdk:"subnet_ids"`
	// SecurityGroupIDs are the security group IDs associated with the image builder.
	SecurityGroupIDs types.Set `tfsdk:"security_group_ids"`
}

type domainJoinInfoModel struct {
	// DirectoryName is the fully qualified name of the Active Directory.
	DirectoryName types.String `tfsdk:"directory_name"`
	// OrganizationalUnitDistinguishedName is the OU DN for computer accounts.
	OrganizationalUnitDistinguishedName types.String `tfsdk:"organizational_unit_distinguished_name"`
}

type accessEndpointModel struct {
	// EndpointType is the type of interface endpoint.
	EndpointType types.String `tfsdk:"endpoint_type"`
	// VpceID is the identifier of the interface VPC endpoint.
	VpceID types.String `tfsdk:"vpce_id"`
}

type rootVolumeConfigModel struct {
	// VolumeSizeInGB is the size of the root volume in GiB.
	VolumeSizeInGB types.Int32 `tfsdk:"volume_size_in_gb"`
}

type networkAccessConfigurationModel struct {
	// EniPrivateIPAddress is the private IP address of the network interface.
	EniPrivateIPAddress types.String `tfsdk:"eni_private_ip_address"`
	// EniIPv6Addresses are the IPv6 addresses assigned to the network interface.
	EniIPv6Addresses types.Set `tfsdk:"eni_ipv6_addresses"`
	// EniID is the identifier of the elastic network interface.
	EniID types.String `tfsdk:"eni_id"`
}

type stateChangeReasonModel struct {
	// Code is the state change reason code (computed).
	Code types.String `tfsdk:"code"`
	// Message is the human-readable state change reason message (computed).
	Message types.String `tfsdk:"message"`
}

type imageBuilderErrorModel struct {
	// ErrorCode is the error code reported by AWS (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
	// ErrorTimestamp is the timestamp when the error occurred (computed).
	ErrorTimestamp types.String `tfsdk:"error_timestamp"`
}
