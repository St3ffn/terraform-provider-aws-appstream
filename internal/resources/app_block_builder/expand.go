// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandVPCConfig(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *awstypes.VpcConfig {
	var m vpcConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	vpcConfig := &awstypes.VpcConfig{
		SubnetIds:        util.ExpandStringSetOrNil(ctx, m.SubnetIDs, diags),
		SecurityGroupIds: util.ExpandStringSetOrNil(ctx, m.SecurityGroupIDs, diags),
	}

	if vpcConfig.SubnetIds == nil && vpcConfig.SecurityGroupIds == nil {
		return nil
	}

	return vpcConfig
}

func expandAccessEndpoints(ctx context.Context, set types.Set, diags *diag.Diagnostics) []awstypes.AccessEndpoint {
	var models []accessEndpointModel
	diags.Append(set.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	if len(models) == 0 {
		return nil
	}

	out := make([]awstypes.AccessEndpoint, 0, len(models))
	for _, m := range models {
		out = append(out, awstypes.AccessEndpoint{
			EndpointType: awstypes.AccessEndpointType(m.EndpointType.ValueString()),
			VpceId:       util.StringPointerOrNil(m.VpceID),
		})
	}

	return out
}
