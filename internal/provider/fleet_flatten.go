// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var fleetComputeCapacityObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"desired_instances": types.Int32Type,
		"desired_sessions":  types.Int32Type,
	},
}

var fleetVPCConfigObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"subnet_ids":         types.SetType{ElemType: types.StringType},
		"security_group_ids": types.SetType{ElemType: types.StringType},
	},
}

var fleetDomainJoinInfoObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"directory_name":                         types.StringType,
		"organizational_unit_distinguished_name": types.StringType,
	},
}

var fleetSessionScriptS3LocationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"s3_bucket": types.StringType,
		"s3_key":    types.StringType,
	},
}

var fleetRootVolumeConfigObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"volume_size_in_gb": types.Int32Type,
	},
}

var fleetErrorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenFleetComputeCapacity(
	ctx context.Context, awsComputeCapacityStatus *awstypes.ComputeCapacityStatus, diags *diag.Diagnostics,
) types.Object {

	if awsComputeCapacityStatus == nil {
		return types.ObjectNull(fleetComputeCapacityObjectType.AttrTypes)
	}

	model := fleetComputeCapacityModel{
		DesiredInstances: int32OrNull(awsComputeCapacityStatus.Desired),
		DesiredSessions:  int32OrNull(awsComputeCapacityStatus.DesiredUserSessions),
	}

	if model.DesiredInstances.IsNull() && model.DesiredSessions.IsNull() {
		return types.ObjectNull(fleetComputeCapacityObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		fleetComputeCapacityObjectType.AttrTypes,
		model,
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(fleetComputeCapacityObjectType.AttrTypes)
	}

	return obj
}

func flattenFleetVPCConfig(
	ctx context.Context, awsVPCConfig *awstypes.VpcConfig, diags *diag.Diagnostics,
) types.Object {

	if awsVPCConfig == nil {
		return types.ObjectNull(fleetVPCConfigObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		fleetVPCConfigObjectType.AttrTypes,
		fleetVPCConfigModel{
			SubnetIDs:        setStringOrNull(ctx, awsVPCConfig.SubnetIds, diags),
			SecurityGroupIDs: setStringOrNull(ctx, awsVPCConfig.SecurityGroupIds, diags),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(fleetVPCConfigObjectType.AttrTypes)
	}

	return obj
}

func flattenFleetDomainJoinInfo(
	ctx context.Context, awsDomainJoinInfo *awstypes.DomainJoinInfo, diags *diag.Diagnostics,
) types.Object {

	if awsDomainJoinInfo == nil {
		return types.ObjectNull(fleetDomainJoinInfoObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		fleetDomainJoinInfoObjectType.AttrTypes,
		fleetDomainJoinInfoModel{
			DirectoryName:                       stringOrNull(awsDomainJoinInfo.DirectoryName),
			OrganizationalUnitDistinguishedName: stringOrNull(awsDomainJoinInfo.OrganizationalUnitDistinguishedName),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(fleetDomainJoinInfoObjectType.AttrTypes)
	}

	return obj
}

func flattenFleetSessionScriptS3Location(
	ctx context.Context, awsS3Location *awstypes.S3Location, diags *diag.Diagnostics,
) types.Object {

	if awsS3Location == nil {
		return types.ObjectNull(fleetSessionScriptS3LocationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		fleetSessionScriptS3LocationObjectType.AttrTypes,
		fleetSessionScriptS3LocationModel{
			S3Bucket: stringOrNull(awsS3Location.S3Bucket),
			S3Key:    stringOrNull(awsS3Location.S3Key),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(fleetSessionScriptS3LocationObjectType.AttrTypes)
	}

	return obj
}

func flattenFleetRootVolumeConfig(
	ctx context.Context, awsVolumeConfig *awstypes.VolumeConfig, diags *diag.Diagnostics,
) types.Object {

	if awsVolumeConfig == nil {
		return types.ObjectNull(fleetRootVolumeConfigObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		fleetRootVolumeConfigObjectType.AttrTypes,
		fleetRootVolumeConfigModel{
			VolumeSizeInGB: int32OrNull(awsVolumeConfig.VolumeSizeInGb),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(fleetRootVolumeConfigObjectType.AttrTypes)
	}

	return obj
}

func flattenFleetErrors(
	ctx context.Context, awsFleetErrors []awstypes.FleetError, diags *diag.Diagnostics,
) types.Set {

	if len(awsFleetErrors) == 0 {
		return types.SetNull(fleetErrorObjectType)
	}

	out := make([]fleetErrorModel, 0, len(awsFleetErrors))
	for _, e := range awsFleetErrors {
		out = append(out, fleetErrorModel{
			ErrorCode:    types.StringValue(string(e.ErrorCode)),
			ErrorMessage: stringOrNull(e.ErrorMessage),
		})
	}

	setVal, d := types.SetValueFrom(ctx, fleetErrorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(fleetErrorObjectType)
	}

	return setVal
}
