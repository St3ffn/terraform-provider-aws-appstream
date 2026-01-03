// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestFlattenStateOwnedBool(t *testing.T) {
	tests := []struct {
		name     string
		prior    types.Bool
		awsValue *bool
		want     types.Bool
	}{
		{
			name:     "prior_null_returns_null",
			prior:    types.BoolNull(),
			awsValue: aws.Bool(true),
			want:     types.BoolNull(),
		},
		{
			name:     "prior_unknown_returns_unknown",
			prior:    types.BoolUnknown(),
			awsValue: aws.Bool(true),
			want:     types.BoolUnknown(),
		},
		{
			name:     "prior_owned_aws_nil_returns_null",
			prior:    types.BoolValue(true),
			awsValue: nil,
			want:     types.BoolNull(),
		},
		{
			name:     "prior_owned_aws_value_returns_value",
			prior:    types.BoolValue(false),
			awsValue: aws.Bool(true),
			want:     types.BoolValue(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenStateOwnedBool(tt.prior, tt.awsValue)
			require.True(t, got.Equal(tt.want), "got %v, want %v", got, tt.want)
		})
	}
}

func TestFlattenStateOwnedInt32(t *testing.T) {
	tests := []struct {
		name     string
		prior    types.Int32
		awsValue *int32
		want     types.Int32
	}{
		{
			name:     "prior_null_returns_null",
			prior:    types.Int32Null(),
			awsValue: aws.Int32(10),
			want:     types.Int32Null(),
		},
		{
			name:     "prior_unknown_returns_unknown",
			prior:    types.Int32Unknown(),
			awsValue: aws.Int32(10),
			want:     types.Int32Unknown(),
		},
		{
			name:     "prior_owned_aws_nil_returns_null",
			prior:    types.Int32Value(5),
			awsValue: nil,
			want:     types.Int32Null(),
		},
		{
			name:     "prior_owned_aws_value_returns_value",
			prior:    types.Int32Value(5),
			awsValue: aws.Int32(42),
			want:     types.Int32Value(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenStateOwnedInt32(tt.prior, tt.awsValue)
			require.True(t, got.Equal(tt.want), "got %v, want %v", got, tt.want)
		})
	}
}

func TestFlattenStateOwnedString(t *testing.T) {
	tests := []struct {
		name     string
		prior    types.String
		awsValue *string
		want     types.String
	}{
		{
			name:     "prior_null_returns_null",
			prior:    types.StringNull(),
			awsValue: aws.String("aws"),
			want:     types.StringNull(),
		},
		{
			name:     "prior_unknown_returns_unknown",
			prior:    types.StringUnknown(),
			awsValue: aws.String("aws"),
			want:     types.StringUnknown(),
		},
		{
			name:     "prior_owned_aws_nil_returns_null",
			prior:    types.StringValue("user"),
			awsValue: nil,
			want:     types.StringNull(),
		},
		{
			name:     "prior_owned_aws_value_returns_value",
			prior:    types.StringValue("user"),
			awsValue: aws.String("aws"),
			want:     types.StringValue("aws"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenStateOwnedString(tt.prior, tt.awsValue)
			require.True(t, got.Equal(tt.want), "got %v, want %v", got, tt.want)
		})
	}
}

func TestFlattenStateOwnedStringSet(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prior     types.Set
		awsValues []string
		want      types.Set
	}{
		{
			name:      "prior_null_returns_null",
			prior:     types.SetNull(types.StringType),
			awsValues: []string{"a"},
			want:      types.SetNull(types.StringType),
		},
		{
			name:      "prior_unknown_returns_unknown",
			prior:     types.SetUnknown(types.StringType),
			awsValues: []string{"a"},
			want:      types.SetUnknown(types.StringType),
		},
		{
			name:      "prior_owned_aws_nil_returns_empty_set",
			prior:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("x")}),
			awsValues: nil,
			want:      types.SetValueMust(types.StringType, []attr.Value{}),
		},
		{
			name:      "prior_owned_aws_empty_returns_empty_set",
			prior:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("x")}),
			awsValues: []string{},
			want:      types.SetValueMust(types.StringType, []attr.Value{}),
		},
		{
			name:      "prior_owned_single_value",
			prior:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("x")}),
			awsValues: []string{"one"},
			want: types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("one")},
			),
		},
		{
			name:      "prior_owned_multiple_values",
			prior:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("x")}),
			awsValues: []string{"one", "two"},
			want: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("one"),
					types.StringValue("two"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := FlattenStateOwnedStringSet(ctx, tt.prior, tt.awsValues, &diags)

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
			require.True(t, got.Equal(tt.want), "got %v, want %v", got, tt.want)
		})
	}
}
