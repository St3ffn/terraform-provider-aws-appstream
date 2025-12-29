// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandComputeCapacity(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.ComputeCapacity {
	var m computeCapacityModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	computeCapacity := &awstypes.ComputeCapacity{
		DesiredInstances: util.Int32PointerOrNil(m.DesiredInstances),
		DesiredSessions:  util.Int32PointerOrNil(m.DesiredSessions),
	}

	if computeCapacity.DesiredInstances == nil && computeCapacity.DesiredSessions == nil {
		return nil
	}

	return computeCapacity
}

func expandVPCConfig(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.VpcConfig {
	var m vpcConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.VpcConfig{
		SubnetIds:        util.ExpandStringSetOrNil(ctx, m.SubnetIDs, diags),
		SecurityGroupIds: util.ExpandStringSetOrNil(ctx, m.SecurityGroupIDs, diags),
	}
}

func expandDomainJoinInfo(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.DomainJoinInfo {
	var m domainJoinInfoModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.DomainJoinInfo{
		DirectoryName:                       util.StringPointerOrNil(m.DirectoryName),
		OrganizationalUnitDistinguishedName: util.StringPointerOrNil(m.OrganizationalUnitDistinguishedName),
	}
}

func expandSessionScriptS3Location(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.S3Location {
	var m sessionScriptS3LocationModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.S3Location{
		S3Bucket: util.StringPointerOrNil(m.S3Bucket),
		S3Key:    util.StringPointerOrNil(m.S3Key),
	}
}

func expandRootVolumeConfig(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.VolumeConfig {
	var m rootVolumeConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	volumeConfig := &awstypes.VolumeConfig{
		VolumeSizeInGb: util.Int32PointerOrNil(m.VolumeSizeInGB),
	}

	if volumeConfig.VolumeSizeInGb == nil {
		return nil
	}

	return volumeConfig
}
