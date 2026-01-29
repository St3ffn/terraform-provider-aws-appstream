// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var associationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
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

func flattenAssociations(
	ctx context.Context,
	awsAssociations []awstypes.SoftwareAssociations,
	diags *diag.Diagnostics,
) types.Map {

	if len(awsAssociations) == 0 {
		return types.MapNull(associationObjectType)
	}

	out := make(map[string]associationModel)

	for _, assoc := range awsAssociations {
		if assoc.SoftwareName == nil {
			continue
		}

		softwareName := aws.ToString(assoc.SoftwareName)

		// deployment errors
		errs := make([]deploymentErrorModel, 0, len(assoc.DeploymentError))
		for _, e := range assoc.DeploymentError {
			errs = append(errs, deploymentErrorModel{
				ErrorCode:    util.StringOrNull(e.ErrorCode),
				ErrorMessage: util.StringOrNull(e.ErrorMessage),
			})
		}

		errSet, d := types.SetValueFrom(ctx, deploymentErrorObjectType, errs)
		diags.Append(d...)
		if diags.HasError() {
			continue
		}

		out[softwareName] = associationModel{
			Status:           types.StringValue(string(assoc.Status)),
			DeploymentErrors: errSet,
		}
	}

	mapVal, d := types.MapValueFrom(ctx, associationObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.MapNull(associationObjectType)
	}

	return mapVal
}
