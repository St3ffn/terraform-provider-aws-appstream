// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package sessions

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenSessions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	expirationTime := time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		in   []awstypes.Session
		want types.Set
	}{
		{
			name: "empty_input",
			in:   nil,
			want: types.SetNull(sessionObjectType),
		},
		{
			name: "single_session_minimal",
			in: []awstypes.Session{
				{
					Id:        aws.String("session-1"),
					UserId:    aws.String("user-1"),
					StackName: aws.String("stack-1"),
					FleetName: aws.String("fleet-1"),
					State:     awstypes.SessionStateActive,
				},
			},
			want: types.SetValueMust(
				sessionObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						sessionObjectType.AttrTypes,
						map[string]attr.Value{
							"id":                  types.StringValue("session-1"),
							"user_id":             types.StringValue("user-1"),
							"stack_name":          types.StringValue("stack-1"),
							"fleet_name":          types.StringValue("fleet-1"),
							"state":               types.StringValue("ACTIVE"),
							"connection_state":    types.StringValue(""),
							"start_time":          types.StringNull(),
							"max_expiration_time": types.StringNull(),
							"authentication_type": types.StringValue(""),
							"instance_id":         types.StringNull(),
							"network_access_configuration": types.ObjectNull(
								networkAccessConfigurationObjectType.AttrTypes,
							),
						},
					),
				},
			),
		},
		{
			name: "single_session_full",
			in: []awstypes.Session{
				{
					Id:                 aws.String("session-2"),
					UserId:             aws.String("user-2"),
					StackName:          aws.String("stack-2"),
					FleetName:          aws.String("fleet-2"),
					State:              awstypes.SessionStateActive,
					ConnectionState:    awstypes.SessionConnectionStateConnected,
					AuthenticationType: awstypes.AuthenticationTypeApi,
					InstanceId:         aws.String("i-1234567890"),
					StartTime:          aws.Time(startTime),
					MaxExpirationTime:  aws.Time(expirationTime),
					NetworkAccessConfiguration: &awstypes.NetworkAccessConfiguration{
						EniId:               aws.String("eni-123"),
						EniPrivateIpAddress: aws.String("10.0.0.10"),
						EniIpv6Addresses:    []string{"2001:db8::1"},
					},
				},
			},
			want: types.SetValueMust(
				sessionObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						sessionObjectType.AttrTypes,
						map[string]attr.Value{
							"id":                  types.StringValue("session-2"),
							"user_id":             types.StringValue("user-2"),
							"stack_name":          types.StringValue("stack-2"),
							"fleet_name":          types.StringValue("fleet-2"),
							"state":               types.StringValue("ACTIVE"),
							"connection_state":    types.StringValue("CONNECTED"),
							"authentication_type": types.StringValue("API"),
							"instance_id":         types.StringValue("i-1234567890"),
							"start_time":          types.StringValue(startTime.Format(time.RFC3339)),
							"max_expiration_time": types.StringValue(expirationTime.Format(time.RFC3339)),
							"network_access_configuration": types.ObjectValueMust(
								networkAccessConfigurationObjectType.AttrTypes,
								map[string]attr.Value{
									"eni_id":                 types.StringValue("eni-123"),
									"eni_private_ip_address": types.StringValue("10.0.0.10"),
									"eni_ipv6_addresses": types.SetValueMust(
										types.StringType,
										[]attr.Value{
											types.StringValue("2001:db8::1"),
										},
									),
								},
							),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenSessions(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenNetworkAccessConfiguration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.NetworkAccessConfiguration
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(networkAccessConfigurationObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.NetworkAccessConfiguration{
				EniId:               aws.String("eni-abc"),
				EniPrivateIpAddress: aws.String("10.0.0.5"),
				EniIpv6Addresses:    []string{"2001:db8::2"},
			},
			want: types.ObjectValueMust(
				networkAccessConfigurationObjectType.AttrTypes,
				map[string]attr.Value{
					"eni_id":                 types.StringValue("eni-abc"),
					"eni_private_ip_address": types.StringValue("10.0.0.5"),
					"eni_ipv6_addresses": types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("2001:db8::2"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenNetworkAccessConfiguration(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
