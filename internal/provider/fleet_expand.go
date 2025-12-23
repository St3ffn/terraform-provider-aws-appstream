// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func expandFleetComputeCapacity(
	ctx context.Context,
	obj types.Object,
	diags *diag.Diagnostics,
) *awstypes.ComputeCapacity {

	var m fleetComputeCapacityModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	computeCapacity := &awstypes.ComputeCapacity{
		DesiredInstances: int32PointerOrNil(m.DesiredInstances),
		DesiredSessions:  int32PointerOrNil(m.DesiredSessions),
	}

	if computeCapacity.DesiredInstances == nil && computeCapacity.DesiredSessions == nil {
		return nil
	}

	return computeCapacity
}

func expandFleetVPCConfig(
	ctx context.Context,
	obj types.Object,
	diags *diag.Diagnostics,
) *awstypes.VpcConfig {

	var m fleetVPCConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.VpcConfig{
		SubnetIds:        expandStringSetOrNil(ctx, m.SubnetIDs, diags),
		SecurityGroupIds: expandStringSetOrNil(ctx, m.SecurityGroupIDs, diags),
	}
}

func expandFleetDomainJoinInfo(
	ctx context.Context,
	obj types.Object,
	diags *diag.Diagnostics,
) *awstypes.DomainJoinInfo {

	var m fleetDomainJoinInfoModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.DomainJoinInfo{
		DirectoryName:                       stringPointerOrNil(m.DirectoryName),
		OrganizationalUnitDistinguishedName: stringPointerOrNil(m.OrganizationalUnitDistinguishedName),
	}
}

func expandFleetSessionScriptS3Location(
	ctx context.Context,
	obj types.Object,
	diags *diag.Diagnostics,
) *awstypes.S3Location {

	var m fleetSessionScriptS3LocationModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.S3Location{
		S3Bucket: stringPointerOrNil(m.S3Bucket),
		S3Key:    stringPointerOrNil(m.S3Key),
	}
}

func expandFleetRootVolumeConfig(
	ctx context.Context,
	obj types.Object,
	diags *diag.Diagnostics,
) *awstypes.VolumeConfig {

	var m fleetRootVolumeConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	volumeConfig := &awstypes.VolumeConfig{
		VolumeSizeInGb: int32PointerOrNil(m.VolumeSizeInGB),
	}

	if volumeConfig.VolumeSizeInGb == nil {
		return nil
	}

	return volumeConfig
}
