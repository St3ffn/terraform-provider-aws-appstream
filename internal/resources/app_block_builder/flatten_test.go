// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

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
	t.Parallel()

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
			t.Parallel()

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

func TestFlattenAccessEndpoints(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

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

func TestFlattenStateChangeReason(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.AppBlockBuilderStateChangeReason
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(stateChangeReasonObjectType.AttrTypes),
		},
		{
			name: "values_set",
			in: &awstypes.AppBlockBuilderStateChangeReason{
				Code:    awstypes.AppBlockBuilderStateChangeReasonCodeInternalError,
				Message: aws.String("boom"),
			},
			want: types.ObjectValueMust(
				stateChangeReasonObjectType.AttrTypes,
				map[string]attr.Value{
					"code":    types.StringValue(string(awstypes.AppBlockBuilderStateChangeReasonCodeInternalError)),
					"message": types.StringValue("boom"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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

func TestFlattenAppBlockBuilderErrors(t *testing.T) {
	t.Parallel()

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
			want: types.SetNull(appBlockBuilderErrorObjectType),
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
				appBlockBuilderErrorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						appBlockBuilderErrorObjectType.AttrTypes,
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
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenAppBlockBuilderErrors(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
