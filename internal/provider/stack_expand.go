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

func expandStorageConnectors(
	ctx context.Context,
	setVal types.Set,
	diags *diag.Diagnostics,
) []awstypes.StorageConnector {

	var models []stackStorageConnectorModel
	diags.Append(setVal.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	out := make([]awstypes.StorageConnector, 0, len(models))
	for _, m := range models {
		out = append(out, awstypes.StorageConnector{
			ConnectorType:              awstypes.StorageConnectorType(m.ConnectorType.ValueString()),
			ResourceIdentifier:         stringPointerOrNil(m.ResourceIdentifier),
			Domains:                    expandStringSetOrNil(ctx, m.Domains, diags),
			DomainsRequireAdminConsent: expandStringSetOrNil(ctx, m.DomainsRequireAdminConsent, diags),
		})
	}

	return out
}

func expandUserSettings(
	ctx context.Context,
	setVal types.Set,
	diags *diag.Diagnostics,
) []awstypes.UserSetting {

	var models []stackUserSettingModel
	diags.Append(setVal.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	out := make([]awstypes.UserSetting, 0, len(models))
	for _, m := range models {
		s := awstypes.UserSetting{
			Action:        awstypes.Action(m.Action.ValueString()),
			Permission:    awstypes.Permission(m.Permission.ValueString()),
			MaximumLength: int32PointerOrNil(m.MaximumLength),
		}

		out = append(out, s)
	}

	return out
}

func expandApplicationSettings(
	ctx context.Context,
	obj types.Object,
	diags *diag.Diagnostics,
) *awstypes.ApplicationSettings {

	var m stackApplicationSettingsModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.ApplicationSettings{
		Enabled:       boolPointerOrNil(m.Enabled),
		SettingsGroup: stringPointerOrNil(m.SettingsGroup),
	}
}

func expandAccessEndpoints(
	ctx context.Context,
	setVal types.Set,
	diags *diag.Diagnostics,
) []awstypes.AccessEndpoint {

	var models []stackAccessEndpointModel
	diags.Append(setVal.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	out := make([]awstypes.AccessEndpoint, 0, len(models))
	for _, m := range models {
		out = append(out, awstypes.AccessEndpoint{
			EndpointType: awstypes.AccessEndpointType(m.EndpointType.ValueString()),
			VpceId:       stringPointerOrNil(m.VpceID),
		})
	}

	return out
}

func expandStreamingExperienceSettings(
	ctx context.Context,
	obj types.Object,
	diags *diag.Diagnostics,
) *awstypes.StreamingExperienceSettings {

	var m stackStreamingExperienceSettingsModel
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
