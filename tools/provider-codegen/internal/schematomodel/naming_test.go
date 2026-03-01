// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import "testing"

func TestToGoFieldName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: "image_arn", want: "ImageARN"},
		{in: "vpc_config", want: "VPCConfig"},
		{in: "max_user_duration_in_seconds", want: "MaxUserDurationInSeconds"},
		{in: "s3_bucket", want: "S3Bucket"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got := toGoFieldName(tt.in)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNestedModelTypeName(t *testing.T) {
	t.Parallel()

	got := nestedModelTypeName("dataSourceModel", []string{
		"post_setup_script_details",
		"script_s3_location",
	})
	want := "dataSourceModelPostSetupScriptDetailsScriptS3Location"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLowerFirst(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: "VPCConfig", want: "vpcConfig"},
		{in: "ARN", want: "arn"},
		{in: "Name", want: "name"},
		{in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got := lowerFirst(tt.in)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUpperFirst(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: "resourceDiff", want: "ResourceDiff"},
		{in: "x", want: "X"},
		{in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got := upperFirst(tt.in)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
