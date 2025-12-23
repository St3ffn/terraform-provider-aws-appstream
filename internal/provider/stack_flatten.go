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

var stackStorageConnectorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"connector_type":                types.StringType,
		"resource_identifier":           types.StringType,
		"domains":                       types.SetType{ElemType: types.StringType},
		"domains_require_admin_consent": types.SetType{ElemType: types.StringType},
	},
}

var stackUserSettingObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"action":         types.StringType,
		"permission":     types.StringType,
		"maximum_length": types.Int32Type,
	},
}

var stackApplicationSettingsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"enabled":        types.BoolType,
		"settings_group": types.StringType,
		"s3_bucket_name": types.StringType,
	},
}

var stackAccessEndpointObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"endpoint_type": types.StringType,
		"vpce_id":       types.StringType,
	},
}

var stackStreamingExperienceSettingsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"preferred_protocol": types.StringType,
	},
}

var stackErrorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenStorageConnectors(
	ctx context.Context, awsStorageConnectors []awstypes.StorageConnector, diags *diag.Diagnostics,
) types.Set {

	if len(awsStorageConnectors) == 0 {
		return types.SetNull(stackStorageConnectorObjectType)
	}

	out := make([]stackStorageConnectorModel, 0, len(awsStorageConnectors))

	for _, c := range awsStorageConnectors {
		m := stackStorageConnectorModel{
			ConnectorType:              types.StringValue(string(c.ConnectorType)),
			ResourceIdentifier:         stringOrNull(c.ResourceIdentifier),
			Domains:                    setStringOrNull(ctx, c.Domains, diags),
			DomainsRequireAdminConsent: setStringOrNull(ctx, c.DomainsRequireAdminConsent, diags),
		}

		out = append(out, m)
	}

	setVal, d := types.SetValueFrom(ctx, stackStorageConnectorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(stackStorageConnectorObjectType)
	}

	return setVal
}

func flattenUserSettings(
	ctx context.Context, awsUserSettings []awstypes.UserSetting, diags *diag.Diagnostics,
) types.Set {

	if len(awsUserSettings) == 0 {
		return types.SetNull(stackUserSettingObjectType)
	}

	out := make([]stackUserSettingModel, 0, len(awsUserSettings))

	for _, s := range awsUserSettings {
		m := stackUserSettingModel{
			Action:        types.StringValue(string(s.Action)),
			Permission:    types.StringValue(string(s.Permission)),
			MaximumLength: int32OrNull(s.MaximumLength),
		}
		out = append(out, m)
	}

	setVal, d := types.SetValueFrom(ctx, stackUserSettingObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(stackUserSettingObjectType)
	}

	return setVal
}

func flattenApplicationSettings(
	ctx context.Context, awsApplicationSettingsResponse *awstypes.ApplicationSettingsResponse, diags *diag.Diagnostics,
) types.Object {

	if awsApplicationSettingsResponse == nil {
		return types.ObjectNull(stackApplicationSettingsObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		stackApplicationSettingsObjectType.AttrTypes,
		stackApplicationSettingsModel{
			Enabled:       boolOrNull(awsApplicationSettingsResponse.Enabled),
			SettingsGroup: stringOrNull(awsApplicationSettingsResponse.SettingsGroup),
			S3BucketName:  stringOrNull(awsApplicationSettingsResponse.S3BucketName),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(stackApplicationSettingsObjectType.AttrTypes)
	}

	return obj
}

func flattenAccessEndpoints(
	ctx context.Context, awsAccessEndpoint []awstypes.AccessEndpoint, diags *diag.Diagnostics,
) types.Set {

	if len(awsAccessEndpoint) == 0 {
		return types.SetNull(stackAccessEndpointObjectType)
	}

	out := make([]stackAccessEndpointModel, 0, len(awsAccessEndpoint))

	for _, e := range awsAccessEndpoint {
		out = append(out, stackAccessEndpointModel{
			EndpointType: types.StringValue(string(e.EndpointType)),
			VpceID:       stringOrNull(e.VpceId),
		})
	}

	setVal, d := types.SetValueFrom(ctx, stackAccessEndpointObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(stackAccessEndpointObjectType)
	}

	return setVal
}

func flattenStreamingExperienceSettings(
	ctx context.Context, awsStreamingExperienceSettings *awstypes.StreamingExperienceSettings, diags *diag.Diagnostics,
) types.Object {

	if awsStreamingExperienceSettings == nil {
		return types.ObjectNull(stackStreamingExperienceSettingsObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		stackStreamingExperienceSettingsObjectType.AttrTypes,
		stackStreamingExperienceSettingsModel{
			PreferredProtocol: types.StringValue(string(awsStreamingExperienceSettings.PreferredProtocol)),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(stackStreamingExperienceSettingsObjectType.AttrTypes)
	}

	return obj
}

func flattenStackErrors(
	ctx context.Context, awsStackErrors []awstypes.StackError, diags *diag.Diagnostics,
) types.Set {

	if len(awsStackErrors) == 0 {
		return types.SetNull(stackErrorObjectType)
	}

	out := make([]stackErrorModel, 0, len(awsStackErrors))

	for _, e := range awsStackErrors {
		out = append(out, stackErrorModel{
			ErrorCode:    types.StringValue(string(e.ErrorCode)),
			ErrorMessage: stringOrNull(e.ErrorMessage),
		})
	}

	setVal, d := types.SetValueFrom(ctx, stackErrorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(stackErrorObjectType)
	}

	return setVal
}
