// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

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

var domainJoinInfoObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"directory_name":                         types.StringType,
		"organizational_unit_distinguished_name": types.StringType,
	},
}

var accessEndpointObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"endpoint_type": types.StringType,
		"vpce_id":       types.StringType,
	},
}

var rootVolumeConfigObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"volume_size_in_gb": types.Int32Type,
	},
}

var networkAccessConfigurationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"eni_private_ip_address": types.StringType,
		"eni_ipv6_addresses":     types.SetType{ElemType: types.StringType},
		"eni_id":                 types.StringType,
	},
}

var stateChangeReasonObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"code":    types.StringType,
		"message": types.StringType,
	},
}

var imageBuilderErrorObjectType = types.ObjectType{
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

func flattenDomainJoinInfo(
	ctx context.Context, awsDomainJoinInfo *awstypes.DomainJoinInfo, diags *diag.Diagnostics,
) types.Object {

	if awsDomainJoinInfo == nil {
		return types.ObjectNull(domainJoinInfoObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		domainJoinInfoObjectType.AttrTypes,
		domainJoinInfoModel{
			DirectoryName:                       util.StringOrNull(awsDomainJoinInfo.DirectoryName),
			OrganizationalUnitDistinguishedName: util.StringOrNull(awsDomainJoinInfo.OrganizationalUnitDistinguishedName),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(domainJoinInfoObjectType.AttrTypes)
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

func flattenRootVolumeConfig(
	ctx context.Context, awsVolumeConfig *awstypes.VolumeConfig, diags *diag.Diagnostics,
) types.Object {

	if awsVolumeConfig == nil {
		return types.ObjectNull(rootVolumeConfigObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		rootVolumeConfigObjectType.AttrTypes,
		rootVolumeConfigModel{
			VolumeSizeInGB: util.Int32OrNull(awsVolumeConfig.VolumeSizeInGb),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(rootVolumeConfigObjectType.AttrTypes)
	}

	return obj
}

func flattenNetworkAccessConfiguration(
	ctx context.Context, awsConfig *awstypes.NetworkAccessConfiguration, diags *diag.Diagnostics,
) types.Object {

	if awsConfig == nil {
		return types.ObjectNull(networkAccessConfigurationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		networkAccessConfigurationObjectType.AttrTypes,
		networkAccessConfigurationModel{
			EniPrivateIPAddress: util.StringOrNull(awsConfig.EniPrivateIpAddress),
			EniIPv6Addresses:    util.SetStringOrNull(ctx, awsConfig.EniIpv6Addresses, diags),
			EniID:               util.StringOrNull(awsConfig.EniId),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(networkAccessConfigurationObjectType.AttrTypes)
	}

	return obj
}

func flattenStateChangeReason(
	ctx context.Context, awsReason *awstypes.ImageBuilderStateChangeReason, diags *diag.Diagnostics,
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

func flattenImageBuilderErrors(
	ctx context.Context, awsErrors []awstypes.ResourceError, diags *diag.Diagnostics,
) types.Set {

	if len(awsErrors) == 0 {
		return types.SetNull(imageBuilderErrorObjectType)
	}

	out := make([]imageBuilderErrorModel, 0, len(awsErrors))
	for _, e := range awsErrors {
		out = append(out, imageBuilderErrorModel{
			ErrorCode:      types.StringValue(string(e.ErrorCode)),
			ErrorMessage:   util.StringOrNull(e.ErrorMessage),
			ErrorTimestamp: util.StringFromTime(e.ErrorTimestamp),
		})
	}

	setVal, d := types.SetValueFrom(ctx, imageBuilderErrorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(imageBuilderErrorObjectType)
	}

	return setVal
}
