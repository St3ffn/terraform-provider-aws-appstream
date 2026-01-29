// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// ID is a synthetic identifier composed of "<name>".
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream app block builder (required).
	Name types.String `tfsdk:"name"`
	// InstanceType is the instance type used to launch the app block builder (computed).
	InstanceType types.String `tfsdk:"instance_type"`
	// Platform is the operating system platform of the app block builder (computed).
	Platform types.String `tfsdk:"platform"`
	// Description is a description to display for the app block builder (computed).
	Description types.String `tfsdk:"description"`
	// DisplayName is the name of the app block builder shown to users (computed).
	DisplayName types.String `tfsdk:"display_name"`
	// VPCConfig specifies the VPC configuration for the app block builder (computed).
	VPCConfig types.Object `tfsdk:"vpc_config"`
	// IAMRoleARN is the ARN of the IAM role applied to the app block builder (computed).
	IAMRoleARN types.String `tfsdk:"iam_role_arn"`
	// EnableDefaultInternetAccess specifies whether the app block builder has internet access (computed).
	EnableDefaultInternetAccess types.Bool `tfsdk:"enable_default_internet_access"`
	// AccessEndpoints specifies interface VPC endpoints used to access the app block builder (computed).
	AccessEndpoints types.Set `tfsdk:"access_endpoints"`
	// Tags is a map of tags assigned to the app block builder (computed).
	Tags types.Map `tfsdk:"tags"`
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
