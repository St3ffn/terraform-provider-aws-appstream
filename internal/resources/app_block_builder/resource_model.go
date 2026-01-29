// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import "github.com/hashicorp/terraform-plugin-framework/types"

type resourceModel struct {
	// ID is a synthetic identifier composed of "<name>".
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream app block builder (required).
	Name types.String `tfsdk:"name"`
	// InstanceType is the instance type used to launch the app block builder (required).
	InstanceType types.String `tfsdk:"instance_type"`
	// Platform is the operating system platform of the app block builder (required).
	Platform types.String `tfsdk:"platform"`
	// Description is a description to display for the app block builder (optional).
	Description types.String `tfsdk:"description"`
	// DisplayName is the name of the app block builder shown to users (optional).
	DisplayName types.String `tfsdk:"display_name"`
	// VPCConfig specifies the VPC configuration for the app block builder (required).
	VPCConfig types.Object `tfsdk:"vpc_config"`
	// IAMRoleARN is the ARN of the IAM role applied to the app block builder (optional).
	IAMRoleARN types.String `tfsdk:"iam_role_arn"`
	// EnableDefaultInternetAccess specifies whether the app block builder has internet access (optional, computed).
	EnableDefaultInternetAccess types.Bool `tfsdk:"enable_default_internet_access"`
	// AccessEndpoints specifies interface VPC endpoints used to access the app block builder (optional).
	AccessEndpoints types.Set `tfsdk:"access_endpoints"`
	// Tags is a map of tags assigned to the app block builder (optional).
	Tags types.Map `tfsdk:"tags"`
	// TagsAll is a map of tags, including default tags, assigned to the app block builder (computed).
	TagsAll types.Map `tfsdk:"tags_all"`
	// ARN is the ARN of the AppStream app block builder (computed).
	ARN types.String `tfsdk:"arn"`
	// CreatedTime is the timestamp when the app block builder was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// State is the current state of the app block builder (computed).
	State types.String `tfsdk:"state"`
	// StateChangeReason describes the most recent state change, if any (computed).
	StateChangeReason types.Object `tfsdk:"state_change_reason"`
	// AppBlockBuilderErrors is the list of errors reported by AWS for the app block builder (computed).
	AppBlockBuilderErrors types.Set `tfsdk:"app_block_builder_errors"`
}

type vpcConfigModel struct {
	// SubnetIDs are the subnet IDs in which the app block builder is launched.
	// App block builders requires subnets in multiple Availability Zones. The constraint is enforced by AWS (optional).
	SubnetIDs types.Set `tfsdk:"subnet_ids"`
	// SecurityGroupIDs are the security group IDs associated with the app block builder (optional).
	SecurityGroupIDs types.Set `tfsdk:"security_group_ids"`
}

type accessEndpointModel struct {
	// EndpointType is the type of interface endpoint (required).
	EndpointType types.String `tfsdk:"endpoint_type"`
	// VpceID is the identifier of the interface VPC endpoint (optional).
	VpceID types.String `tfsdk:"vpce_id"`
}

type stateChangeReasonModel struct {
	// Code is the state change reason code (computed).
	Code types.String `tfsdk:"code"`
	// Message is the human-readable state change reason message (computed).
	Message types.String `tfsdk:"message"`
}

type appBlockBuilderErrorModel struct {
	// ErrorCode is the error code reported by AWS (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
	// ErrorTimestamp is the timestamp when the error occurred (computed).
	ErrorTimestamp types.String `tfsdk:"error_timestamp"`
}
