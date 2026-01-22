// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var vpcConfigObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"subnet_ids":         types.SetType{ElemType: types.StringType},
		"security_group_ids": types.SetType{ElemType: types.StringType},
	},
}

var accessEndpointObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"endpoint_type": types.StringType,
		"vpce_id":       types.StringType,
	},
}

var stateChangeReasonObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"code":    types.StringType,
		"message": types.StringType,
	},
}

var appBlockBuilderErrorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":      types.StringType,
		"error_message":   types.StringType,
		"error_timestamp": types.StringType,
	},
}

func flattenVPCConfig(ctx context.Context, awsVPCConfig *awstypes.VpcConfig, diags *diag.Diagnostics) types.Object {
	if awsVPCConfig == nil {
		return types.ObjectNull(vpcConfigObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		vpcConfigObjectType.AttrTypes,
		vpcConfigModel{
			SubnetIDs:        util.SetStringOrNull(ctx, awsVPCConfig.SubnetIds, diags),
			SecurityGroupIDs: util.SetStringOrNull(ctx, awsVPCConfig.SecurityGroupIds, diags),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(vpcConfigObjectType.AttrTypes)
	}

	return obj
}

func flattenAccessEndpoints(
	ctx context.Context, awsEndpoints []awstypes.AccessEndpoint, diags *diag.Diagnostics,
) types.Set {

	if len(awsEndpoints) == 0 {
		return types.SetNull(accessEndpointObjectType)
	}

	out := make([]accessEndpointModel, 0, len(awsEndpoints))
	for _, e := range awsEndpoints {
		out = append(out, accessEndpointModel{
			EndpointType: types.StringValue(string(e.EndpointType)),
			VpceID:       util.StringOrNull(e.VpceId),
		})
	}

	setVal, d := types.SetValueFrom(ctx, accessEndpointObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(accessEndpointObjectType)
	}

	return setVal
}

func flattenStateChangeReason(
	ctx context.Context, awsReason *awstypes.AppBlockBuilderStateChangeReason, diags *diag.Diagnostics,
) types.Object {

	if awsReason == nil {
		return types.ObjectNull(stateChangeReasonObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		stateChangeReasonObjectType.AttrTypes,
		stateChangeReasonModel{
			Code:    types.StringValue(string(awsReason.Code)),
			Message: util.StringOrNull(awsReason.Message),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(stateChangeReasonObjectType.AttrTypes)
	}

	return obj
}

func flattenAppBlockBuilderErrors(
	ctx context.Context, awsErrors []awstypes.ResourceError, diags *diag.Diagnostics,
) types.Set {

	if len(awsErrors) == 0 {
		return types.SetNull(appBlockBuilderErrorObjectType)
	}

	out := make([]appBlockBuilderErrorModel, 0, len(awsErrors))
	for _, e := range awsErrors {
		out = append(out, appBlockBuilderErrorModel{
			ErrorCode:      types.StringValue(string(e.ErrorCode)),
			ErrorMessage:   util.StringOrNull(e.ErrorMessage),
			ErrorTimestamp: util.StringFromTime(e.ErrorTimestamp),
		})
	}

	setVal, d := types.SetValueFrom(ctx, appBlockBuilderErrorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(appBlockBuilderErrorObjectType)
	}

	return setVal
}
