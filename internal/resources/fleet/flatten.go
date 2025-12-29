// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var computeCapacityObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"desired_instances": types.Int32Type,
		"desired_sessions":  types.Int32Type,
	},
}

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

var sessionScriptS3LocationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"s3_bucket": types.StringType,
		"s3_key":    types.StringType,
	},
}

var rootVolumeConfigObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"volume_size_in_gb": types.Int32Type,
	},
}

var errorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenComputeCapacity(
	ctx context.Context, awsComputeCapacityStatus *awstypes.ComputeCapacityStatus, diags *diag.Diagnostics,
) types.Object {

	if awsComputeCapacityStatus == nil {
		return types.ObjectNull(computeCapacityObjectType.AttrTypes)
	}

	model := computeCapacityModel{
		DesiredInstances: util.Int32OrNull(awsComputeCapacityStatus.Desired),
		DesiredSessions:  util.Int32OrNull(awsComputeCapacityStatus.DesiredUserSessions),
	}

	if model.DesiredInstances.IsNull() && model.DesiredSessions.IsNull() {
		return types.ObjectNull(computeCapacityObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		computeCapacityObjectType.AttrTypes,
		model,
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(computeCapacityObjectType.AttrTypes)
	}

	return obj
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

func flattenSessionScriptS3Location(
	ctx context.Context, awsS3Location *awstypes.S3Location, diags *diag.Diagnostics,
) types.Object {

	if awsS3Location == nil {
		return types.ObjectNull(sessionScriptS3LocationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		sessionScriptS3LocationObjectType.AttrTypes,
		sessionScriptS3LocationModel{
			S3Bucket: util.StringOrNull(awsS3Location.S3Bucket),
			S3Key:    util.StringOrNull(awsS3Location.S3Key),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(sessionScriptS3LocationObjectType.AttrTypes)
	}

	return obj
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

func flattenFleetErrors(ctx context.Context, awsFleetErrors []awstypes.FleetError, diags *diag.Diagnostics) types.Set {
	if len(awsFleetErrors) == 0 {
		return types.SetNull(errorObjectType)
	}

	out := make([]fleetErrorModel, 0, len(awsFleetErrors))
	for _, e := range awsFleetErrors {
		out = append(out, fleetErrorModel{
			ErrorCode:    types.StringValue(string(e.ErrorCode)),
			ErrorMessage: util.StringOrNull(e.ErrorMessage),
		})
	}

	setVal, d := types.SetValueFrom(ctx, errorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(errorObjectType)
	}

	return setVal
}
