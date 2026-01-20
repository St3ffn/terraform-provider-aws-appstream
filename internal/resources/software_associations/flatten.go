// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package software_associations

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var softwareAssociationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"software_name":     types.StringType,
		"status":            types.StringType,
		"deployment_errors": types.SetType{ElemType: deploymentErrorObjectType},
	},
}

var deploymentErrorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenSoftwareAssociations(
	ctx context.Context, awsAssociations []awstypes.SoftwareAssociations, diags *diag.Diagnostics,
) types.Set {

	if len(awsAssociations) == 0 {
		return types.SetNull(softwareAssociationObjectType)
	}

	out := make([]softwareAssociationModel, 0, len(awsAssociations))
	for _, a := range awsAssociations {
		out = append(out, softwareAssociationModel{
			SoftwareName:     util.StringOrNull(a.SoftwareName),
			Status:           types.StringValue(string(a.Status)),
			DeploymentErrors: flattenSoftwareAssociationDeploymentErrors(ctx, a.DeploymentError, diags),
		})
	}

	setVal, d := types.SetValueFrom(ctx, softwareAssociationObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(softwareAssociationObjectType)
	}

	return setVal
}

func flattenSoftwareAssociationDeploymentErrors(
	ctx context.Context, awsErrors []awstypes.ErrorDetails, diags *diag.Diagnostics,
) types.Set {

	if len(awsErrors) == 0 {
		return types.SetNull(deploymentErrorObjectType)
	}

	out := make([]deploymentErrorModel, 0, len(awsErrors))
	for _, e := range awsErrors {
		out = append(out, deploymentErrorModel{
			ErrorCode:    util.StringOrNull(e.ErrorCode),
			ErrorMessage: util.StringOrNull(e.ErrorMessage),
		})
	}

	setVal, d := types.SetValueFrom(ctx, deploymentErrorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(deploymentErrorObjectType)
	}

	return setVal
}
