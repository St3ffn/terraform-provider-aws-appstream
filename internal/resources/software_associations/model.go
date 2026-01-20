// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package software_associations

import "github.com/hashicorp/terraform-plugin-framework/types"

type model struct {
	// AssociatedResource is the ARN of the AppStream image or image builder (required).
	AssociatedResource types.String `tfsdk:"associated_resource"`
	// SoftwareAssociations is the list of license-included software associations
	// for the specified AppStream resource (computed).
	SoftwareAssociations types.Set `tfsdk:"software_associations"`
}

type softwareAssociationModel struct {
	// SoftwareName is the name of the license-included application (computed).
	SoftwareName types.String `tfsdk:"software_name"`
	// Status is the deployment status of the software association (computed).
	Status types.String `tfsdk:"status"`
	// DeploymentErrors is the list of deployment errors reported by AWS, if any (computed).
	DeploymentErrors types.Set `tfsdk:"deployment_errors"`
}

type deploymentErrorModel struct {
	// ErrorCode is the error code reported by AWS for the software deployment (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message reported by AWS (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
}
