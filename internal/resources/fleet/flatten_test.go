// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenComputeCapacity(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.ComputeCapacityStatus
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(computeCapacityObjectType.AttrTypes),
		},
		{
			name: "both_values_set",
			in: &awstypes.ComputeCapacityStatus{
				Desired:             aws.Int32(2),
				DesiredUserSessions: aws.Int32(10),
			},
			want: types.ObjectValueMust(
				computeCapacityObjectType.AttrTypes,
				map[string]attr.Value{
					"desired_instances": types.Int32Value(2),
					"desired_sessions":  types.Int32Value(10),
				},
			),
		},
		{
			name: "both_values_null",
			in:   &awstypes.ComputeCapacityStatus{},
			want: types.ObjectNull(computeCapacityObjectType.AttrTypes),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenComputeCapacity(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestFlattenSessionScriptS3Location(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.S3Location
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(sessionScriptS3LocationObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.S3Location{
				S3Bucket: aws.String("bucket"),
				S3Key:    aws.String("script.ps1"),
			},
			want: types.ObjectValueMust(
				sessionScriptS3LocationObjectType.AttrTypes,
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

			got := flattenSessionScriptS3Location(ctx, tt.in, &diags)

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
				VolumeSizeInGb: aws.Int32(100),
			},
			want: types.ObjectValueMust(
				rootVolumeConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"volume_size_in_gb": types.Int32Value(100),
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

func TestFlattenErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.FleetError
		want types.Set
	}{
		{
			name: "empty_slice",
			in:   nil,
			want: types.SetNull(errorObjectType),
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
				errorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						errorObjectType.AttrTypes,
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
