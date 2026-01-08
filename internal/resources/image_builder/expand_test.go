// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

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

func TestExpandVPCConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		obj  types.Object
		want *awstypes.VpcConfig
	}{
		{
			name: "both_fields_set",
			obj: types.ObjectValueMust(
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
			want: &awstypes.VpcConfig{
				SubnetIds:        []string{"subnet-1"},
				SecurityGroupIds: []string{"sg-1"},
			},
		},
		{
			name: "only_subnets_set",
			obj: types.ObjectValueMust(
				vpcConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"subnet_ids": types.SetValueMust(
						types.StringType,
						[]attr.Value{types.StringValue("subnet-1")},
					),
					"security_group_ids": types.SetNull(types.StringType),
				},
			),
			want: &awstypes.VpcConfig{
				SubnetIds: []string{"subnet-1"},
			},
		},
		{
			name: "both_null_returns_nil",
			obj: types.ObjectValueMust(
				vpcConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"subnet_ids":         types.SetNull(types.StringType),
					"security_group_ids": types.SetNull(types.StringType),
				},
			),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandVPCConfig(ctx, tt.obj, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandDomainJoinInfo(t *testing.T) {
	ctx := context.Background()

	obj := types.ObjectValueMust(
		domainJoinInfoObjectType.AttrTypes,
		map[string]attr.Value{
			"directory_name":                         types.StringValue("corp.example.com"),
			"organizational_unit_distinguished_name": types.StringValue("OU=Apps,DC=example,DC=com"),
		},
	)

	var diags diag.Diagnostics
	got := expandDomainJoinInfo(ctx, obj, &diags)

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if got == nil {
		t.Fatalf("expected DomainJoinInfo, got nil")
	}

	if aws.ToString(got.DirectoryName) != "corp.example.com" {
		t.Fatalf("unexpected directory name: %q", aws.ToString(got.DirectoryName))
	}

	if aws.ToString(got.OrganizationalUnitDistinguishedName) != "OU=Apps,DC=example,DC=com" {
		t.Fatalf("unexpected OU DN")
	}
}

func TestExpandAccessEndpoints(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		set  types.Set
		want []awstypes.AccessEndpoint
	}{
		{
			name: "single_endpoint",
			set: types.SetValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"endpoint_type": types.StringType,
						"vpce_id":       types.StringType,
					},
				},
				[]attr.Value{
					types.ObjectValueMust(
						map[string]attr.Type{
							"endpoint_type": types.StringType,
							"vpce_id":       types.StringType,
						},
						map[string]attr.Value{
							"endpoint_type": types.StringValue("STREAMING"),
							"vpce_id":       types.StringValue("vpce-123"),
						},
					),
				},
			),
			want: []awstypes.AccessEndpoint{
				{
					EndpointType: awstypes.AccessEndpointTypeStreaming,
					VpceId:       aws.String("vpce-123"),
				},
			},
		},
		{
			name: "empty_set_returns_nil",
			set:  types.SetNull(types.ObjectType{}),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandAccessEndpoints(ctx, tt.set, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandRootVolumeConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		obj  types.Object
		want *awstypes.VolumeConfig
	}{
		{
			name: "volume_size_set",
			obj: types.ObjectValueMust(
				rootVolumeConfigObjectType.AttrTypes,
				map[string]attr.Value{
					"volume_size_in_gb": types.Int32Value(200),
				},
			),
			want: &awstypes.VolumeConfig{
				VolumeSizeInGb: aws.Int32(200),
			},
		},
		{
			name: "volume_size_null_returns_nil",
			obj: types.ObjectValueMust(
				rootVolumeConfigObjectType.AttrTypes,
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

			got := expandRootVolumeConfig(ctx, tt.obj, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
