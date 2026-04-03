// Copyright St3ffn 2025, 2026
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

	var models []resourceModelStorageConnectors
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
	var models []resourceModelUserSettings
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

	var m resourceModelApplicationSettings
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.ApplicationSettings{
		Enabled:       util.BoolPointerOrNil(m.Enabled),
		SettingsGroup: util.StringPointerOrNil(m.SettingsGroup),
	}
}

func expandContentRedirection(
	ctx context.Context, obj types.Object, diags *diag.Diagnostics,
) *awstypes.ContentRedirection {

	var m resourceModelContentRedirection
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	if m.HostToClient.IsNull() {
		return nil
	}

	var hostToClient resourceModelContentRedirectionHostToClient
	diags.Append(m.HostToClient.As(ctx, &hostToClient, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.ContentRedirection{
		HostToClient: &awstypes.UrlRedirectionConfig{
			Enabled:     util.BoolPointerOrNil(hostToClient.Enabled),
			AllowedUrls: util.ExpandStringSetOrNil(ctx, hostToClient.AllowedUrls, diags),
			DeniedUrls:  util.ExpandStringSetOrNil(ctx, hostToClient.DeniedUrls, diags),
		},
	}
}

func expandAccessEndpoints(ctx context.Context, setVal types.Set, diags *diag.Diagnostics) []awstypes.AccessEndpoint {
	var models []resourceModelAccessEndpoints
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

	var m resourceModelStreamingExperienceSettings
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
