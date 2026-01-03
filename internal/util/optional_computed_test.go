// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestFlattenOptionalComputedBool(t *testing.T) {
	tests := []struct {
		name     string
		prior    types.Bool
		awsValue *bool
		want     types.Bool
	}{
		{
			name:     "aws_value_overrides_prior_null",
			prior:    types.BoolNull(),
			awsValue: aws.Bool(true),
			want:     types.BoolValue(true),
		},
		{
			name:     "aws_value_overrides_prior_unknown",
			prior:    types.BoolUnknown(),
			awsValue: aws.Bool(true),
			want:     types.BoolValue(true),
		},
		{
			name:     "aws_value_overrides_prior_owned",
			prior:    types.BoolValue(false),
			awsValue: aws.Bool(true),
			want:     types.BoolValue(true),
		},
		{
			name:     "aws_nil_prior_null_returns_null",
			prior:    types.BoolNull(),
			awsValue: nil,
			want:     types.BoolNull(),
		},
		{
			name:     "aws_nil_prior_unknown_returns_unknown",
			prior:    types.BoolUnknown(),
			awsValue: nil,
			want:     types.BoolUnknown(),
		},
		{
			name:     "aws_nil_prior_owned_keeps_prior",
			prior:    types.BoolValue(true),
			awsValue: nil,
			want:     types.BoolValue(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenOptionalComputedBool(tt.prior, tt.awsValue)
			require.True(t, got.Equal(tt.want), "got %v, want %v", got, tt.want)
		})
	}
}

func TestFlattenOptionalComputedString(t *testing.T) {
	tests := []struct {
		name     string
		prior    types.String
		awsValue *string
		want     types.String
	}{
		{
			name:     "aws_value_overrides_prior_null",
			prior:    types.StringNull(),
			awsValue: aws.String("test"),
			want:     types.StringValue("test"),
		},
		{
			name:     "aws_value_overrides_prior_unknown",
			prior:    types.StringUnknown(),
			awsValue: aws.String("test"),
			want:     types.StringValue("test"),
		},
		{
			name:     "aws_value_overrides_prior_owned",
			prior:    types.StringValue("other"),
			awsValue: aws.String("test"),
			want:     types.StringValue("test"),
		},
		{
			name:     "aws_nil_prior_null_returns_null",
			prior:    types.StringNull(),
			awsValue: nil,
			want:     types.StringNull(),
		},
		{
			name:     "aws_nil_prior_unknown_returns_unknown",
			prior:    types.StringUnknown(),
			awsValue: nil,
			want:     types.StringUnknown(),
		},
		{
			name:     "aws_nil_prior_owned_keeps_prior",
			prior:    types.StringValue("test"),
			awsValue: nil,
			want:     types.StringValue("test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenOptionalComputedString(tt.prior, tt.awsValue)
			require.True(t, got.Equal(tt.want), "got %v, want %v", got, tt.want)
		})
	}
}
