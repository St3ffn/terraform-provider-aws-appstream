// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package sessions

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var sessionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                  types.StringType,
		"user_id":             types.StringType,
		"stack_name":          types.StringType,
		"fleet_name":          types.StringType,
		"state":               types.StringType,
		"connection_state":    types.StringType,
		"start_time":          types.StringType,
		"max_expiration_time": types.StringType,
		"authentication_type": types.StringType,
		"instance_id":         types.StringType,
		"network_access_configuration": types.ObjectType{
			AttrTypes: networkAccessConfigurationObjectType.AttrTypes,
		},
	},
}

var networkAccessConfigurationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"eni_private_ip_address": types.StringType,
		"eni_ipv6_addresses":     types.SetType{ElemType: types.StringType},
		"eni_id":                 types.StringType,
	},
}

func flattenSessions(
	ctx context.Context, awsSessions []awstypes.Session, diags *diag.Diagnostics,
) types.Set {

	if len(awsSessions) == 0 {
		return types.SetNull(sessionObjectType)
	}

	out := make([]sessionModel, 0, len(awsSessions))

	for _, session := range awsSessions {
		m := sessionModel{
			ID:                 util.StringOrNull(session.Id),
			UserID:             util.StringOrNull(session.UserId),
			StackName:          util.StringOrNull(session.StackName),
			FleetName:          util.StringOrNull(session.FleetName),
			State:              types.StringValue(string(session.State)),
			ConnectionState:    types.StringValue(string(session.ConnectionState)),
			StartTime:          util.StringFromTime(session.StartTime),
			MaxExpirationTime:  util.StringFromTime(session.MaxExpirationTime),
			AuthenticationType: types.StringValue(string(session.AuthenticationType)),
			InstanceID:         util.StringOrNull(session.InstanceId),
			NetworkAccessConfiguration: flattenNetworkAccessConfiguration(
				ctx, session.NetworkAccessConfiguration, diags,
			),
		}

		out = append(out, m)
	}

	setVal, d := types.SetValueFrom(ctx, sessionObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(sessionObjectType)
	}

	return setVal
}

func flattenNetworkAccessConfiguration(
	ctx context.Context, awsNetwork *awstypes.NetworkAccessConfiguration, diags *diag.Diagnostics,
) types.Object {

	if awsNetwork == nil {
		return types.ObjectNull(networkAccessConfigurationObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		networkAccessConfigurationObjectType.AttrTypes,
		networkAccessConfigurationModel{
			EniPrivateIPAddress: util.StringOrNull(awsNetwork.EniPrivateIpAddress),
			EniIPv6Addresses:    util.SetStringOrNull(ctx, awsNetwork.EniIpv6Addresses, diags),
			EniID:               util.StringOrNull(awsNetwork.EniId),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(networkAccessConfigurationObjectType.AttrTypes)
	}

	return obj
}
