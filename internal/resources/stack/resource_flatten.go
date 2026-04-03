// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var storageConnectorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"connector_type":                types.StringType,
		"resource_identifier":           types.StringType,
		"domains":                       types.SetType{ElemType: types.StringType},
		"domains_require_admin_consent": types.SetType{ElemType: types.StringType},
	},
}

var userSettingObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"action":         types.StringType,
		"permission":     types.StringType,
		"maximum_length": types.Int32Type,
	},
}

var applicationSettingsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"enabled":        types.BoolType,
		"settings_group": types.StringType,
		"s3_bucket_name": types.StringType,
	},
}

var contentRedirectionHostToClientObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"enabled":      types.BoolType,
		"allowed_urls": types.SetType{ElemType: types.StringType},
		"denied_urls":  types.SetType{ElemType: types.StringType},
	},
}

var contentRedirectionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"host_to_client": contentRedirectionHostToClientObjectType,
	},
}

var accessEndpointObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"endpoint_type": types.StringType,
		"vpce_id":       types.StringType,
	},
}

var streamingExperienceSettingsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"preferred_protocol": types.StringType,
	},
}

var errorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenStorageConnectorsResource(
	ctx context.Context, prior types.Set, awsConnectors []awstypes.StorageConnector, diags *diag.Diagnostics,
) types.Set {

	// user never managed this attribute
	if prior.IsNull() {
		return types.SetNull(storageConnectorObjectType)
	}

	// terraform does not yet know the value during planning
	if prior.IsUnknown() {
		return types.SetUnknown(storageConnectorObjectType)
	}

	// decode prior state into models
	var priorModels []resourceModelStorageConnectors
	diags.Append(prior.ElementsAs(ctx, &priorModels, false)...)
	if diags.HasError() {
		return types.SetNull(storageConnectorObjectType)
	}

	// index aws connectors by identity using connector_type
	awsByType := make(map[string]awstypes.StorageConnector, len(awsConnectors))
	for _, awsConn := range awsConnectors {
		if awsConn.ConnectorType == "" {
			continue
		}
		awsByType[string(awsConn.ConnectorType)] = awsConn
	}

	out := make([]resourceModelStorageConnectors, 0, len(priorModels))

	// iterate prior state to preserve ownership and intent
	for _, priorConn := range priorModels {
		if priorConn.ConnectorType.IsNull() || priorConn.ConnectorType.IsUnknown() {
			continue
		}

		ct := priorConn.ConnectorType.ValueString()
		awsConn, exists := awsByType[ct]

		// drift case: user configured the connector but aws no longer has it
		if !exists {
			out = append(out, resourceModelStorageConnectors{
				ConnectorType:              types.StringValue(ct),
				ResourceIdentifier:         types.StringNull(),
				Domains:                    types.SetNull(types.StringType),
				DomainsRequireAdminConsent: types.SetNull(types.StringType),
			})
			continue
		}

		m := resourceModelStorageConnectors{
			ConnectorType:      types.StringValue(ct),
			ResourceIdentifier: util.FlattenStateOwnedString(priorConn.ResourceIdentifier, awsConn.ResourceIdentifier),
			Domains:            util.FlattenStateOwnedStringSet(ctx, priorConn.Domains, awsConn.Domains, diags),
			DomainsRequireAdminConsent: util.FlattenStateOwnedStringSet(
				ctx, priorConn.DomainsRequireAdminConsent, awsConn.DomainsRequireAdminConsent, diags,
			),
		}

		out = append(out, m)
	}

	if diags.HasError() {
		return types.SetNull(storageConnectorObjectType)
	}

	// user explicitly wants zero items
	if len(out) == 0 {
		empty, d := types.SetValue(storageConnectorObjectType, []attr.Value{})
		diags.Append(d...)
		return empty
	}

	setVal, d := types.SetValueFrom(ctx, storageConnectorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(storageConnectorObjectType)
	}

	return setVal
}

func flattenUserSettingsResource(
	ctx context.Context, prior types.Set, awsUserSettings []awstypes.UserSetting, diags *diag.Diagnostics,
) types.Set {

	// user never managed this attribute
	if prior.IsNull() {
		return types.SetNull(userSettingObjectType)
	}

	// terraform does not yet know during planning
	if prior.IsUnknown() {
		return types.SetUnknown(userSettingObjectType)
	}

	// decode prior state into models
	var priorModels []resourceModelUserSettings
	diags.Append(prior.ElementsAs(ctx, &priorModels, false)...)
	if diags.HasError() {
		return types.SetNull(userSettingObjectType)
	}

	// index aws user settings by action
	awsByAction := make(map[string]awstypes.UserSetting, len(awsUserSettings))
	for _, s := range awsUserSettings {
		awsByAction[string(s.Action)] = s
	}

	out := make([]resourceModelUserSettings, 0, len(priorModels))

	// iterate prior state to preserve ownership
	for _, priorSetting := range priorModels {
		if priorSetting.Action.IsNull() || priorSetting.Action.IsUnknown() {
			continue
		}

		action := priorSetting.Action.ValueString()
		awsSetting, exists := awsByAction[action]

		// drift: user configured action but aws no longer has it
		if !exists {
			out = append(out, resourceModelUserSettings{
				Action:        types.StringValue(action),
				Permission:    types.StringNull(),
				MaximumLength: types.Int32Null(),
			})
			continue
		}

		out = append(out, resourceModelUserSettings{
			Action:        types.StringValue(action),
			Permission:    types.StringValue(string(awsSetting.Permission)),
			MaximumLength: util.Int32OrNull(awsSetting.MaximumLength),
		})
	}

	// user explicitly wants zero items
	if len(out) == 0 {
		empty, d := types.SetValue(userSettingObjectType, []attr.Value{})
		diags.Append(d...)
		return empty
	}

	setVal, d := types.SetValueFrom(ctx, userSettingObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(userSettingObjectType)
	}

	return setVal
}

func flattenApplicationSettingsResource(
	ctx context.Context, prior types.Object, awsAppSettings *awstypes.ApplicationSettingsResponse, diags *diag.Diagnostics,
) types.Object {

	// user never managed this attribute
	if prior.IsNull() {
		return types.ObjectNull(applicationSettingsObjectType.AttrTypes)
	}

	// terraform does not yet know during planning
	if prior.IsUnknown() {
		return types.ObjectUnknown(applicationSettingsObjectType.AttrTypes)
	}

	// decode prior state
	var priorModel resourceModelApplicationSettings
	diags.Append(prior.As(ctx, &priorModel, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return types.ObjectNull(applicationSettingsObjectType.AttrTypes)
	}

	// drift: user configured it but aws no longer has it
	if awsAppSettings == nil {
		obj, d := types.ObjectValueFrom(
			ctx,
			applicationSettingsObjectType.AttrTypes,
			resourceModelApplicationSettings{
				Enabled:       priorModel.Enabled,
				SettingsGroup: types.StringNull(),
				S3BucketName:  types.StringNull(),
			},
		)
		diags.Append(d...)
		return obj
	}

	// normal reconcile
	obj, d := types.ObjectValueFrom(
		ctx,
		applicationSettingsObjectType.AttrTypes,
		resourceModelApplicationSettings{
			Enabled:       util.BoolOrNull(awsAppSettings.Enabled),
			SettingsGroup: util.FlattenStateOwnedString(priorModel.SettingsGroup, awsAppSettings.SettingsGroup),
			S3BucketName:  util.StringOrNull(awsAppSettings.S3BucketName),
		},
	)
	diags.Append(d...)
	return obj
}

func flattenContentRedirectionResource(
	ctx context.Context, prior types.Object, awsContentRedirection *awstypes.ContentRedirection, diags *diag.Diagnostics,
) types.Object {

	// user never managed this attribute
	if prior.IsNull() {
		return types.ObjectNull(contentRedirectionObjectType.AttrTypes)
	}

	// terraform does not yet know during planning
	if prior.IsUnknown() {
		return types.ObjectUnknown(contentRedirectionObjectType.AttrTypes)
	}

	// decode prior state
	var priorModel resourceModelContentRedirection
	diags.Append(prior.As(ctx, &priorModel, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return types.ObjectNull(contentRedirectionObjectType.AttrTypes)
	}

	// drift: user configured it but aws no longer has it
	if awsContentRedirection == nil || awsContentRedirection.HostToClient == nil {
		obj, d := types.ObjectValueFrom(
			ctx,
			contentRedirectionObjectType.AttrTypes,
			resourceModelContentRedirection{
				HostToClient: types.ObjectNull(contentRedirectionHostToClientObjectType.AttrTypes),
			},
		)
		diags.Append(d...)
		return obj
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		contentRedirectionObjectType.AttrTypes,
		resourceModelContentRedirection{
			HostToClient: mustObjectValueFromStack(
				ctx,
				contentRedirectionHostToClientObjectType.AttrTypes,
				resourceModelContentRedirectionHostToClient{
					Enabled: util.BoolOrNull(awsContentRedirection.HostToClient.Enabled),
					AllowedUrls: util.FlattenStateOwnedStringSet(
						ctx,
						extractContentRedirectionHostToClient(ctx, priorModel).AllowedUrls,
						awsContentRedirection.HostToClient.AllowedUrls, diags,
					),
					DeniedUrls: util.FlattenStateOwnedStringSet(
						ctx,
						extractContentRedirectionHostToClient(ctx, priorModel).DeniedUrls,
						awsContentRedirection.HostToClient.DeniedUrls, diags,
					),
				},
				diags,
			),
		},
	)
	diags.Append(d...)
	return obj
}

func extractContentRedirectionHostToClient(
	ctx context.Context, prior resourceModelContentRedirection,
) resourceModelContentRedirectionHostToClient {

	if prior.HostToClient.IsNull() || prior.HostToClient.IsUnknown() {
		return resourceModelContentRedirectionHostToClient{
			Enabled:     types.BoolNull(),
			AllowedUrls: types.SetNull(types.StringType),
			DeniedUrls:  types.SetNull(types.StringType),
		}
	}

	var out resourceModelContentRedirectionHostToClient
	_ = prior.HostToClient.As(ctx, &out, basetypes.ObjectAsOptions{})
	return out
}

func mustObjectValueFromStack[T any](
	ctx context.Context, attrTypes map[string]attr.Type, in T, diags *diag.Diagnostics,
) types.Object {
	obj, d := types.ObjectValueFrom(ctx, attrTypes, in)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(attrTypes)
	}
	return obj
}

func flattenAccessEndpointsResource(
	ctx context.Context, prior types.Set, awsAccessEndpoints []awstypes.AccessEndpoint, diags *diag.Diagnostics,
) types.Set {

	// user never managed this attribute
	if prior.IsNull() {
		return types.SetNull(accessEndpointObjectType)
	}

	// terraform does not yet know during planning
	if prior.IsUnknown() {
		return types.SetUnknown(accessEndpointObjectType)
	}

	// aws returned nothing but user manages it -> drift
	if len(awsAccessEndpoints) == 0 {
		empty, d := types.SetValue(accessEndpointObjectType, []attr.Value{})
		diags.Append(d...)
		return empty
	}

	out := make([]resourceModelAccessEndpoints, 0, len(awsAccessEndpoints))

	for _, e := range awsAccessEndpoints {
		out = append(out, resourceModelAccessEndpoints{
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

func flattenStreamingExperienceSettingsResource(
	ctx context.Context, prior types.Object, awsSettings *awstypes.StreamingExperienceSettings, diags *diag.Diagnostics,
) types.Object {

	// user never managed this attribute
	if prior.IsNull() {
		return types.ObjectNull(streamingExperienceSettingsObjectType.AttrTypes)
	}

	// terraform does not yet know during planning
	if prior.IsUnknown() {
		return types.ObjectUnknown(streamingExperienceSettingsObjectType.AttrTypes)
	}

	// decode prior state
	var priorModel resourceModelStreamingExperienceSettings
	diags.Append(prior.As(ctx, &priorModel, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return types.ObjectNull(streamingExperienceSettingsObjectType.AttrTypes)
	}

	// drift: user configured it but aws no longer has it
	if awsSettings == nil {
		obj, d := types.ObjectValueFrom(
			ctx,
			streamingExperienceSettingsObjectType.AttrTypes,
			resourceModelStreamingExperienceSettings{
				PreferredProtocol: priorModel.PreferredProtocol,
			},
		)
		diags.Append(d...)
		return obj
	}

	// normal reconcile (ignore aws default unless user owns block)
	obj, d := types.ObjectValueFrom(
		ctx,
		streamingExperienceSettingsObjectType.AttrTypes,
		resourceModelStreamingExperienceSettings{
			PreferredProtocol: types.StringValue(string(awsSettings.PreferredProtocol)),
		},
	)
	diags.Append(d...)
	return obj
}
