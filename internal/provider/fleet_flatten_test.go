// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenFleetComputeCapacity(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.ComputeCapacityStatus
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(fleetComputeCapacityObjectType.AttrTypes),
		},
		{
			name: "both_values_set",
			in: &awstypes.ComputeCapacityStatus{
				Desired:             aws.Int32(2),
				DesiredUserSessions: aws.Int32(10),
			},
			want: types.ObjectValueMust(
				fleetComputeCapacityObjectType.AttrTypes,
				map[string]attr.Value{
					"desired_instances": types.Int32Value(2),
					"desired_sessions":  types.Int32Value(10),
				},
			),
		},
		{
			name: "both_values_null",
			in:   &awstypes.ComputeCapacityStatus{},
			want: types.ObjectNull(fleetComputeCapacityObjectType.AttrTypes),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenFleetComputeCapacity(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenFleetVPCConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.VpcConfig
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(fleetVPCConfigObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.VpcConfig{
				SubnetIds:        []string{"subnet-1"},
				SecurityGroupIds: []string{"sg-1"},
			},
			want: types.ObjectValueMust(
				fleetVPCConfigObjectType.AttrTypes,
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

			got := flattenFleetVPCConfig(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenFleetDomainJoinInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.DomainJoinInfo
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(fleetDomainJoinInfoObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.DomainJoinInfo{
				DirectoryName:                       aws.String("example.com"),
				OrganizationalUnitDistinguishedName: aws.String("OU=Apps"),
			},
			want: types.ObjectValueMust(
				fleetDomainJoinInfoObjectType.AttrTypes,
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

			got := flattenFleetDomainJoinInfo(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenFleetSessionScriptS3Location(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.S3Location
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(fleetSessionScriptS3LocationObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.S3Location{
				S3Bucket: aws.String("bucket"),
				S3Key:    aws.String("script.ps1"),
			},
			want: types.ObjectValueMust(
				fleetSessionScriptS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("bucket"),
					"s3_key":    types.StringValue("script.ps1"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenFleetSessionScriptS3Location(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenFleetRootVolumeConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.VolumeConfig
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(fleetRootVolumeConfigObjectType.AttrTypes),
		},
		{
			name: "volume_size_set",
			in: &awstypes.VolumeConfig{
				VolumeSizeInGb: aws.Int32(100),
			},
			want: types.ObjectValueMust(
				fleetRootVolumeConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"volume_size_in_gb": types.Int32Value(100),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenFleetRootVolumeConfig(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenFleetErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.FleetError
		want types.Set
	}{
		{
			name: "empty_slice",
			in:   nil,
			want: types.SetNull(fleetErrorObjectType),
		},
		{
			name: "single_error",
			in: []awstypes.FleetError{
				{
					ErrorCode:    awstypes.FleetErrorCodeInternalServiceError,
					ErrorMessage: aws.String("boom"),
				},
			},
			want: types.SetValueMust(
				fleetErrorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						fleetErrorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":    types.StringValue("INTERNAL_SERVICE_ERROR"),
							"error_message": types.StringValue("boom"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenFleetErrors(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
