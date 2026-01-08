// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

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

func TestFlattenVPCConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.VpcConfig
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(vpcConfigObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.VpcConfig{
				SubnetIds:        []string{"subnet-1"},
				SecurityGroupIds: []string{"sg-1"},
			},
			want: types.ObjectValueMust(
				vpcConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"subnet_ids": types.SetValueMust(
						types.StringType,
						[]attr.Value{types.StringValue("subnet-1")},
					),
					"security_group_ids": types.SetValueMust(
						types.StringType,
						[]attr.Value{types.StringValue("sg-1")},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenVPCConfig(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenDomainJoinInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.DomainJoinInfo
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(domainJoinInfoObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.DomainJoinInfo{
				DirectoryName:                       aws.String("example.com"),
				OrganizationalUnitDistinguishedName: aws.String("OU=Apps"),
			},
			want: types.ObjectValueMust(
				domainJoinInfoObjectType.AttrTypes,
				map[string]attr.Value{
					"directory_name":                         types.StringValue("example.com"),
					"organizational_unit_distinguished_name": types.StringValue("OU=Apps"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenDomainJoinInfo(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenAccessEndpoints(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.AccessEndpoint
		want types.Set
	}{
		{
			name: "empty_slice",
			in:   nil,
			want: types.SetNull(accessEndpointObjectType),
		},
		{
			name: "single_endpoint",
			in: []awstypes.AccessEndpoint{
				{
					EndpointType: awstypes.AccessEndpointTypeStreaming,
					VpceId:       aws.String("vpce-123"),
				},
			},
			want: types.SetValueMust(
				accessEndpointObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						accessEndpointObjectType.AttrTypes,
						map[string]attr.Value{
							"endpoint_type": types.StringValue("STREAMING"),
							"vpce_id":       types.StringValue("vpce-123"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenAccessEndpoints(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenRootVolumeConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.VolumeConfig
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(rootVolumeConfigObjectType.AttrTypes),
		},
		{
			name: "volume_size_set",
			in: &awstypes.VolumeConfig{
				VolumeSizeInGb: aws.Int32(200),
			},
			want: types.ObjectValueMust(
				rootVolumeConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"volume_size_in_gb": types.Int32Value(200),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenRootVolumeConfig(ctx, tt.in, &diags)

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
				EniPrivateIpAddress: aws.String("10.0.0.5"),
				EniIpv6Addresses:    []string{"2001:db8::1"},
				EniId:               aws.String("eni-123"),
			},
			want: types.ObjectValueMust(
				networkAccessConfigurationObjectType.AttrTypes,
				map[string]attr.Value{
					"eni_private_ip_address": types.StringValue("10.0.0.5"),
					"eni_ipv6_addresses": types.SetValueMust(
						types.StringType,
						[]attr.Value{types.StringValue("2001:db8::1")},
					),
					"eni_id": types.StringValue("eni-123"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestFlattenStateChangeReason(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.ImageBuilderStateChangeReason
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(stateChangeReasonObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.ImageBuilderStateChangeReason{
				Code:    awstypes.ImageBuilderStateChangeReasonCodeInternalError,
				Message: aws.String("boom"),
			},
			want: types.ObjectValueMust(
				stateChangeReasonObjectType.AttrTypes,
				map[string]attr.Value{
					"code":    types.StringValue(string(awstypes.ImageBuilderStateChangeReasonCodeInternalError)),
					"message": types.StringValue("boom"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenStateChangeReason(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenImageBuilderErrors(t *testing.T) {
	ctx := context.Background()
	ts := time.Now()

	tests := []struct {
		name string
		in   []awstypes.ResourceError
		want types.Set
	}{
		{
			name: "empty_slice",
			in:   nil,
			want: types.SetNull(imageBuilderErrorObjectType),
		},
		{
			name: "single_error",
			in: []awstypes.ResourceError{
				{
					ErrorCode:      awstypes.FleetErrorCodeInternalServiceError,
					ErrorMessage:   aws.String("boom"),
					ErrorTimestamp: aws.Time(ts),
				},
			},
			want: types.SetValueMust(
				imageBuilderErrorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						imageBuilderErrorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":      types.StringValue(string(awstypes.FleetErrorCodeInternalServiceError)),
							"error_message":   types.StringValue("boom"),
							"error_timestamp": types.StringValue(ts.Format(time.RFC3339)),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenImageBuilderErrors(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
