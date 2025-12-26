// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpandFleetComputeCapacity(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		obj  types.Object
		want *awstypes.ComputeCapacity
	}{
		{
			name: "both_fields_set",
			obj: types.ObjectValueMust(
				fleetComputeCapacityObjectType.AttrTypes,
				map[string]attr.Value{
					"desired_instances": types.Int32Value(2),
					"desired_sessions":  types.Int32Value(10),
				},
			),
			want: &awstypes.ComputeCapacity{
				DesiredInstances: aws.Int32(2),
				DesiredSessions:  aws.Int32(10),
			},
		},
		{
			name: "only_instances_set",
			obj: types.ObjectValueMust(
				fleetComputeCapacityObjectType.AttrTypes,
				map[string]attr.Value{
					"desired_instances": types.Int32Value(1),
					"desired_sessions":  types.Int32Null(),
				},
			),
			want: &awstypes.ComputeCapacity{
				DesiredInstances: aws.Int32(1),
			},
		},
		{
			name: "both_null_returns_nil",
			obj: types.ObjectValueMust(
				fleetComputeCapacityObjectType.AttrTypes,
				map[string]attr.Value{
					"desired_instances": types.Int32Null(),
					"desired_sessions":  types.Int32Null(),
				},
			),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandFleetComputeCapacity(ctx, tt.obj, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandFleetVPCConfig(t *testing.T) {
	ctx := context.Background()

	obj := types.ObjectValueMust(
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
	)

	var diags diag.Diagnostics
	got := expandFleetVPCConfig(ctx, obj, &diags)

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if got == nil {
		t.Fatalf("expected VpcConfig, got nil")
	}

	if !reflect.DeepEqual(got.SubnetIds, []string{"subnet-1"}) {
		t.Fatalf("unexpected subnet ids: %#v", got.SubnetIds)
	}

	if !reflect.DeepEqual(got.SecurityGroupIds, []string{"sg-1"}) {
		t.Fatalf("unexpected security group ids: %#v", got.SecurityGroupIds)
	}
}

func TestExpandFleetDomainJoinInfo(t *testing.T) {
	ctx := context.Background()

	obj := types.ObjectValueMust(
		fleetDomainJoinInfoObjectType.AttrTypes,
		map[string]attr.Value{
			"directory_name":                         types.StringValue("example.com"),
			"organizational_unit_distinguished_name": types.StringValue("OU=Apps,DC=example,DC=com"),
		},
	)

	var diags diag.Diagnostics
	got := expandFleetDomainJoinInfo(ctx, obj, &diags)

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if got == nil {
		t.Fatalf("expected DomainJoinInfo, got nil")
	}

	if aws.ToString(got.DirectoryName) != "example.com" {
		t.Fatalf("unexpected directory name")
	}
}

func TestExpandFleetSessionScriptS3Location(t *testing.T) {
	ctx := context.Background()

	obj := types.ObjectValueMust(
		fleetSessionScriptS3LocationObjectType.AttrTypes,
		map[string]attr.Value{
			"s3_bucket": types.StringValue("bucket"),
			"s3_key":    types.StringValue("script.ps1"),
		},
	)

	var diags diag.Diagnostics
	got := expandFleetSessionScriptS3Location(ctx, obj, &diags)

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if got == nil {
		t.Fatalf("expected S3Location, got nil")
	}

	if aws.ToString(got.S3Bucket) != "bucket" {
		t.Fatalf("unexpected bucket")
	}
}

func TestExpandFleetRootVolumeConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		obj  types.Object
		want *awstypes.VolumeConfig
	}{
		{
			name: "volume_size_set",
			obj: types.ObjectValueMust(
				fleetRootVolumeConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"volume_size_in_gb": types.Int32Value(100),
				},
			),
			want: &awstypes.VolumeConfig{
				VolumeSizeInGb: aws.Int32(100),
			},
		},
		{
			name: "volume_size_null_returns_nil",
			obj: types.ObjectValueMust(
				fleetRootVolumeConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"volume_size_in_gb": types.Int32Null(),
				},
			),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandFleetRootVolumeConfig(ctx, tt.obj, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
