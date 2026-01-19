// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import "github.com/hashicorp/terraform-plugin-framework/types"

type resourceModel struct {
	// ID is a synthetic identifier composed of "<image_builder_arn>" (computed).
	ID types.String `tfsdk:"id"`
	// ImageBuilderARN is the ARN of the AppStream image builder to associate software with (required).
	ImageBuilderARN types.String `tfsdk:"image_builder_arn"`
	// SoftwareNames is the set of license-included software package names
	// associated with the image builder (required).
	SoftwareNames types.Set `tfsdk:"software_names"`
	// Deploy controls whether a software deployment is triggered after association
	// (optional, computed; defaults to false).
	Deploy types.Bool `tfsdk:"deploy"`
	// Associations contains per-software association status and deployment details
	// as reported by AWS (computed).
	Associations types.Map `tfsdk:"associations"`
}

type associationModel struct {
	// Status is the AWS-reported deployment status of the license-included
	// application for the image builder (computed).
	Status types.String `tfsdk:"status"`
	// DeploymentErrors contains error details reported by AWS for failed
	// software deployments (computed).
	DeploymentErrors types.Set `tfsdk:"deployment_errors"`
}

type deploymentErrorModel struct {
	// ErrorCode is the AWS-reported error code (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable AWS error message (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
}
