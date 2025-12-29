// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandStorageConnectors(
	ctx context.Context, setVal types.Set, diags *diag.Diagnostics,
) []awstypes.StorageConnector {

	var models []storageConnectorModel
	diags.Append(setVal.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	out := make([]awstypes.StorageConnector, 0, len(models))
	for _, m := range models {
		out = append(out, awstypes.StorageConnector{
			ConnectorType:              awstypes.StorageConnectorType(m.ConnectorType.ValueString()),
			ResourceIdentifier:         util.StringPointerOrNil(m.ResourceIdentifier),
			Domains:                    util.ExpandStringSetOrNil(ctx, m.Domains, diags),
			DomainsRequireAdminConsent: util.ExpandStringSetOrNil(ctx, m.DomainsRequireAdminConsent, diags),
		})
	}

	return out
}

func expandUserSettings(ctx context.Context, setVal types.Set, diags *diag.Diagnostics) []awstypes.UserSetting {
	var models []userSettingModel
	diags.Append(setVal.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	out := make([]awstypes.UserSetting, 0, len(models))
	for _, m := range models {
		s := awstypes.UserSetting{
			Action:        awstypes.Action(m.Action.ValueString()),
			Permission:    awstypes.Permission(m.Permission.ValueString()),
			MaximumLength: util.Int32PointerOrNil(m.MaximumLength),
		}

		out = append(out, s)
	}

	return out
}

func expandApplicationSettings(
	ctx context.Context, obj types.Object, diags *diag.Diagnostics,
) *awstypes.ApplicationSettings {

	var m applicationSettingsModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.ApplicationSettings{
		Enabled:       util.BoolPointerOrNil(m.Enabled),
		SettingsGroup: util.StringPointerOrNil(m.SettingsGroup),
	}
}

func expandAccessEndpoints(ctx context.Context, setVal types.Set, diags *diag.Diagnostics) []awstypes.AccessEndpoint {
	var models []accessEndpointModel
	diags.Append(setVal.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
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

func expandStreamingExperienceSettings(
	ctx context.Context, obj types.Object, diags *diag.Diagnostics,
) *awstypes.StreamingExperienceSettings {

	var m streamingExperienceSettingsModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	if m.PreferredProtocol.IsNull() {
		return nil
	}

	return &awstypes.StreamingExperienceSettings{
		PreferredProtocol: awstypes.PreferredProtocol(m.PreferredProtocol.ValueString()),
	}
}
