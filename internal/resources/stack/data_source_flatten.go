// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func flattenStorageConnectorsData(
	ctx context.Context, awsStorageConnectors []awstypes.StorageConnector, diags *diag.Diagnostics,
) types.Set {

	if len(awsStorageConnectors) == 0 {
		return types.SetNull(storageConnectorObjectType)
	}

	out := make([]storageConnectorModel, 0, len(awsStorageConnectors))

	for _, c := range awsStorageConnectors {
		m := storageConnectorModel{
			ConnectorType:              types.StringValue(string(c.ConnectorType)),
			ResourceIdentifier:         util.StringOrNull(c.ResourceIdentifier),
			Domains:                    util.SetStringOrNull(ctx, c.Domains, diags),
			DomainsRequireAdminConsent: util.SetStringOrNull(ctx, c.DomainsRequireAdminConsent, diags),
		}

		out = append(out, m)
	}

	setVal, d := types.SetValueFrom(ctx, storageConnectorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(storageConnectorObjectType)
	}

	return setVal
}

func flattenUserSettingsData(
	ctx context.Context, awsUserSettings []awstypes.UserSetting, diags *diag.Diagnostics,
) types.Set {

	if len(awsUserSettings) == 0 {
		return types.SetNull(userSettingObjectType)
	}

	out := make([]userSettingModel, 0, len(awsUserSettings))

	for _, s := range awsUserSettings {
		m := userSettingModel{
			Action:        types.StringValue(string(s.Action)),
			Permission:    types.StringValue(string(s.Permission)),
			MaximumLength: util.Int32OrNull(s.MaximumLength),
		}
		out = append(out, m)
	}

	setVal, d := types.SetValueFrom(ctx, userSettingObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(userSettingObjectType)
	}

	return setVal
}

func flattenApplicationSettingsData(
	ctx context.Context, awsApplicationSettingsResponse *awstypes.ApplicationSettingsResponse, diags *diag.Diagnostics,
) types.Object {

	if awsApplicationSettingsResponse == nil {
		return types.ObjectNull(applicationSettingsObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		applicationSettingsObjectType.AttrTypes,
		applicationSettingsModel{
			Enabled:       util.BoolOrNull(awsApplicationSettingsResponse.Enabled),
			SettingsGroup: util.StringOrNull(awsApplicationSettingsResponse.SettingsGroup),
			S3BucketName:  util.StringOrNull(awsApplicationSettingsResponse.S3BucketName),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(applicationSettingsObjectType.AttrTypes)
	}

	return obj
}

func flattenAccessEndpointsData(
	ctx context.Context, awsAccessEndpoint []awstypes.AccessEndpoint, diags *diag.Diagnostics,
) types.Set {

	if len(awsAccessEndpoint) == 0 {
		return types.SetNull(accessEndpointObjectType)
	}

	out := make([]accessEndpointModel, 0, len(awsAccessEndpoint))

	for _, e := range awsAccessEndpoint {
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

func flattenStreamingExperienceSettingsData(
	ctx context.Context, awsStreamingExperienceSettings *awstypes.StreamingExperienceSettings, diags *diag.Diagnostics,
) types.Object {

	if awsStreamingExperienceSettings == nil {
		return types.ObjectNull(streamingExperienceSettingsObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		streamingExperienceSettingsObjectType.AttrTypes,
		streamingExperienceSettingsModel{
			PreferredProtocol: types.StringValue(string(awsStreamingExperienceSettings.PreferredProtocol)),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(streamingExperienceSettingsObjectType.AttrTypes)
	}

	return obj
}

func flattenStackErrorsData(ctx context.Context, awsStackErrors []awstypes.StackError, diags *diag.Diagnostics) types.Set {
	if len(awsStackErrors) == 0 {
		return types.SetNull(errorObjectType)
	}

	out := make([]stackErrorModel, 0, len(awsStackErrors))

	for _, e := range awsStackErrors {
		out = append(out, stackErrorModel{
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
