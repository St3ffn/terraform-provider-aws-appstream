// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	awstypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ taggingAPI = (*fakeTaggingAPI)(nil)

type fakeTaggingAPI struct {
	GetResourcesFn func(
		ctx context.Context, params *awstaggingapi.GetResourcesInput, optFns ...func(*awstaggingapi.Options),
	) (*awstaggingapi.GetResourcesOutput, error)

	TagResourcesFn func(
		ctx context.Context, params *awstaggingapi.TagResourcesInput, optFns ...func(*awstaggingapi.Options),
	) (*awstaggingapi.TagResourcesOutput, error)

	UntagResourcesFn func(
		ctx context.Context, params *awstaggingapi.UntagResourcesInput, optFns ...func(*awstaggingapi.Options),
	) (*awstaggingapi.UntagResourcesOutput, error)

	GetResourcesCalls   int
	TagResourcesCalls   int
	UntagResourcesCalls int

	LastGetResourcesInput   *awstaggingapi.GetResourcesInput
	LastTagResourcesInput   *awstaggingapi.TagResourcesInput
	LastUntagResourcesInput *awstaggingapi.UntagResourcesInput
}

func newFakeTaggingAPI() *fakeTaggingAPI {
	return &fakeTaggingAPI{}
}

func (f *fakeTaggingAPI) GetResources(
	ctx context.Context,
	params *awstaggingapi.GetResourcesInput,
	optFns ...func(*awstaggingapi.Options),
) (*awstaggingapi.GetResourcesOutput, error) {
	if f.GetResourcesFn == nil {
		panic("GetResources called but not configured")
	}

	f.GetResourcesCalls++
	f.LastGetResourcesInput = params

	return f.GetResourcesFn(ctx, params, optFns...)
}

func (f *fakeTaggingAPI) TagResources(
	ctx context.Context, params *awstaggingapi.TagResourcesInput, optFns ...func(*awstaggingapi.Options),
) (*awstaggingapi.TagResourcesOutput, error) {

	if f.TagResourcesFn == nil {
		panic("TagResources called but not configured")
	}

	f.TagResourcesCalls++
	f.LastTagResourcesInput = params

	return f.TagResourcesFn(ctx, params, optFns...)
}

func (f *fakeTaggingAPI) UntagResources(
	ctx context.Context, params *awstaggingapi.UntagResourcesInput, optFns ...func(*awstaggingapi.Options),
) (*awstaggingapi.UntagResourcesOutput, error) {

	if f.UntagResourcesFn == nil {
		panic("UntagResources called but not configured")
	}

	f.UntagResourcesCalls++
	f.LastUntagResourcesInput = params

	return f.UntagResourcesFn(ctx, params, optFns...)
}

func (f *fakeTaggingAPI) GetResourcesReturns(out *awstaggingapi.GetResourcesOutput) *fakeTaggingAPI {
	f.GetResourcesFn = func(
		context.Context,
		*awstaggingapi.GetResourcesInput,
		...func(*awstaggingapi.Options),
	) (*awstaggingapi.GetResourcesOutput, error) {
		return out, nil
	}
	return f
}

func (f *fakeTaggingAPI) GetResourcesFails(err error) *fakeTaggingAPI {
	f.GetResourcesFn = func(
		context.Context,
		*awstaggingapi.GetResourcesInput,
		...func(*awstaggingapi.Options),
	) (*awstaggingapi.GetResourcesOutput, error) {
		return nil, err
	}
	return f
}

func (f *fakeTaggingAPI) TagResourcesSucceeds() *fakeTaggingAPI {
	f.TagResourcesFn = func(
		context.Context,
		*awstaggingapi.TagResourcesInput,
		...func(*awstaggingapi.Options),
	) (*awstaggingapi.TagResourcesOutput, error) {
		return &awstaggingapi.TagResourcesOutput{}, nil
	}
	return f
}

func (f *fakeTaggingAPI) TagResourcesFails(err error) *fakeTaggingAPI {
	f.TagResourcesFn = func(
		context.Context,
		*awstaggingapi.TagResourcesInput,
		...func(*awstaggingapi.Options),
	) (*awstaggingapi.TagResourcesOutput, error) {
		return nil, err
	}
	return f
}

func (f *fakeTaggingAPI) UntagResourcesSucceeds() *fakeTaggingAPI {
	f.UntagResourcesFn = func(
		context.Context,
		*awstaggingapi.UntagResourcesInput,
		...func(*awstaggingapi.Options),
	) (*awstaggingapi.UntagResourcesOutput, error) {
		return &awstaggingapi.UntagResourcesOutput{}, nil
	}
	return f
}

func (f *fakeTaggingAPI) UntagResourcesFails(err error) *fakeTaggingAPI {
	f.UntagResourcesFn = func(
		context.Context,
		*awstaggingapi.UntagResourcesInput,
		...func(*awstaggingapi.Options),
	) (*awstaggingapi.UntagResourcesOutput, error) {
		return nil, err
	}
	return f
}

func TestTagManager_Read(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		arn         string
		setupClient func(*fakeTaggingAPI)
		want        types.Map
		wantError   bool
	}{
		{
			name: "empty_arn_returns_null",
			arn:  "",
			setupClient: func(f *fakeTaggingAPI) {
				// no calls expected
			},
			want:      types.MapNull(types.StringType),
			wantError: false,
		},
		{
			name: "successful_read_single_tag",
			arn:  "arn:aws:appstream:eu-central-1:123456789012:stack/test",
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesReturns(&awstaggingapi.GetResourcesOutput{
					ResourceTagMappingList: []awstypes.ResourceTagMapping{
						{
							Tags: []awstypes.Tag{
								{Key: aws.String("env"), Value: aws.String("prod")},
							},
						},
					},
				})
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"env": types.StringValue("prod"),
				},
			),
			wantError: false,
		},
		{
			name: "successful_read_multiple_tags",
			arn:  "arn:aws:appstream:eu-central-1:123456789012:stack/test",
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesReturns(&awstaggingapi.GetResourcesOutput{
					ResourceTagMappingList: []awstypes.ResourceTagMapping{
						{
							Tags: []awstypes.Tag{
								{Key: aws.String("env"), Value: aws.String("prod")},
								{Key: aws.String("team"), Value: aws.String("core")},
							},
						},
					},
				})
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"env":  types.StringValue("prod"),
					"team": types.StringValue("core"),
				},
			),
			wantError: false,
		},
		{
			name: "aws_error_returns_diagnostics",
			arn:  "arn:aws:appstream:eu-central-1:123456789012:stack/test",
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesFails(errors.New("boom"))
			},
			want:      types.MapNull(types.StringType),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := newFakeTaggingAPI()
			tt.setupClient(fake)

			tm := newTagManager(fake, nil)

			got, diags := tm.Read(ctx, tt.arn)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error, got none")
				}
				return
			}

			require.Falsef(
				t, diags.HasError(),
				"unexpected diagnostics for arn %q: %v", tt.arn, diags,
			)

			require.Truef(
				t, got.Equal(tt.want),
				"Read(%q) mismatch\nGot:  %#v\nWant: %#v", tt.arn, got, tt.want,
			)
		})
	}
}
func TestTagManager_Apply(t *testing.T) {
	ctx := context.Background()
	arn := "arn:aws:appstream:eu-central-1:123456789012:stack/test"

	tests := []struct {
		name        string
		arn         string
		defaultTags map[string]string
		desired     types.Map
		setupClient func(*fakeTaggingAPI)
		assert      func(t *testing.T, f *fakeTaggingAPI)
		want        types.Map
		wantError   bool
	}{
		{
			name: "empty_arn_noop",
			arn:  "",
			desired: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{"a": types.StringValue("b")},
			),
			setupClient: func(f *fakeTaggingAPI) {},
			assert: func(t *testing.T, f *fakeTaggingAPI) {
				if f.GetResourcesCalls != 0 || f.TagResourcesCalls != 0 || f.UntagResourcesCalls != 0 {
					t.Fatalf("no AWS calls expected for empty ARN")
				}
			},
			want: types.MapNull(types.StringType),
		},
		{
			name:    "desired_unknown_delegates_to_read",
			arn:     arn,
			desired: types.MapUnknown(types.StringType),
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesReturns(&awstaggingapi.GetResourcesOutput{
					ResourceTagMappingList: []awstypes.ResourceTagMapping{
						{
							Tags: []awstypes.Tag{
								{Key: aws.String("env"), Value: aws.String("prod")},
							},
						},
					},
				})
			},
			assert: func(t *testing.T, f *fakeTaggingAPI) {
				if f.GetResourcesCalls != 1 {
					t.Fatalf("expected GetResources to be called once, got %d", f.GetResourcesCalls)
				}
				if f.TagResourcesCalls != 0 || f.UntagResourcesCalls != 0 {
					t.Fatalf("no tag/untag calls expected when desired is unknown")
				}
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{"env": types.StringValue("prod")},
			),
		},
		{
			name: "add_tags_only",
			arn:  arn,
			desired: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{"env": types.StringValue("prod")},
			),
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesReturns(&awstaggingapi.GetResourcesOutput{})
				f.TagResourcesSucceeds()
			},
			assert: func(t *testing.T, f *fakeTaggingAPI) {
				if f.TagResourcesCalls != 1 {
					t.Fatalf("expected TagResources to be called once")
				}
				if f.UntagResourcesCalls != 0 {
					t.Fatalf("did not expect UntagResources to be called")
				}

				if got := f.LastTagResourcesInput.Tags["env"]; got != "prod" {
					t.Fatalf("expected tag env=prod, got %q", got)
				}
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{"env": types.StringValue("prod")},
			),
		},
		{
			name:    "remove_tags_only",
			arn:     arn,
			desired: types.MapNull(types.StringType),
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesReturns(&awstaggingapi.GetResourcesOutput{
					ResourceTagMappingList: []awstypes.ResourceTagMapping{
						{
							Tags: []awstypes.Tag{
								{Key: aws.String("old"), Value: aws.String("value")},
							},
						},
					},
				})
				f.UntagResourcesSucceeds()
			},
			assert: func(t *testing.T, f *fakeTaggingAPI) {
				if f.UntagResourcesCalls != 1 {
					t.Fatalf("expected UntagResources to be called once")
				}
				if len(f.LastUntagResourcesInput.TagKeys) != 1 ||
					f.LastUntagResourcesInput.TagKeys[0] != "old" {
					t.Fatalf("expected to untag key 'old'")
				}
			},
			want: types.MapNull(types.StringType),
		},
		{
			name:    "add_and_remove_tags",
			arn:     arn,
			desired: types.MapValueMust(types.StringType, map[string]attr.Value{"new": types.StringValue("v")}),
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesReturns(&awstaggingapi.GetResourcesOutput{
					ResourceTagMappingList: []awstypes.ResourceTagMapping{
						{
							Tags: []awstypes.Tag{
								{Key: aws.String("old"), Value: aws.String("v")},
							},
						},
					},
				})
				f.UntagResourcesSucceeds()
				f.TagResourcesSucceeds()
			},
			assert: func(t *testing.T, f *fakeTaggingAPI) {
				if f.UntagResourcesCalls != 1 || f.TagResourcesCalls != 1 {
					t.Fatalf("expected both TagResources and UntagResources to be called once")
				}
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{"new": types.StringValue("v")},
			),
		},
		{
			name:    "untag_error_returns_diagnostics",
			arn:     arn,
			desired: types.MapNull(types.StringType),
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesReturns(&awstaggingapi.GetResourcesOutput{
					ResourceTagMappingList: []awstypes.ResourceTagMapping{
						{
							Tags: []awstypes.Tag{
								{Key: aws.String("a"), Value: aws.String("b")},
							},
						},
					},
				})
				f.UntagResourcesFails(errors.New("boom"))
			},
			wantError: true,
		},
		{
			name:    "tag_error_returns_diagnostics",
			arn:     arn,
			desired: types.MapValueMust(types.StringType, map[string]attr.Value{"a": types.StringValue("b")}),
			setupClient: func(f *fakeTaggingAPI) {
				f.GetResourcesReturns(&awstaggingapi.GetResourcesOutput{})
				f.TagResourcesFails(errors.New("boom"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := newFakeTaggingAPI()
			if tt.setupClient != nil {
				tt.setupClient(fake)
			}

			tm := newTagManager(fake, tt.defaultTags)

			got, diags := tm.Apply(ctx, tt.arn, tt.desired)

			if tt.wantError {
				require.Truef(
					t, diags.HasError(),
					"expected diagnostics error, got none",
				)
				return
			}

			require.Falsef(
				t, diags.HasError(),
				"unexpected diagnostics: %v", diags,
			)

			require.Truef(
				t, got.Equal(tt.want),
				"Apply() result mismatch\nGot:  %#v\nWant: %#v", got, tt.want,
			)

			if tt.assert != nil {
				tt.assert(t, fake)
			}
		})
	}
}

func TestFlattenTags(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     map[string]string
		want      types.Map
		wantError bool
	}{
		{
			name:  "nil_map_returns_null",
			input: nil,
			want:  types.MapNull(types.StringType),
		},
		{
			name:  "empty_map_returns_null",
			input: map[string]string{},
			want:  types.MapNull(types.StringType),
		},
		{
			name: "single_tag",
			input: map[string]string{
				"env": "prod",
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"env": types.StringValue("prod"),
				},
			),
		},
		{
			name: "multiple_tags",
			input: map[string]string{
				"env":  "prod",
				"team": "core",
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"env":  types.StringValue("prod"),
					"team": types.StringValue("core"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenTags(ctx, tt.input, &diags)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error, got none")
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			require.Truef(
				t, got.Equal(tt.want),
				"flattenTags(%v) mismatch\nGot:  %#v\nWant: %#v", tt.input, got, tt.want,
			)
		})
	}
}

func TestExpandTags(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Map
		want      map[string]string
		wantError bool
	}{
		{
			name:  "null_map_returns_nil",
			input: types.MapNull(types.StringType),
			want:  nil,
		},
		{
			name:  "empty_map_returns_empty",
			input: types.MapValueMust(types.StringType, map[string]attr.Value{}),
			want:  map[string]string{},
		},
		{
			name: "single_tag",
			input: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"env": types.StringValue("prod"),
				},
			),
			want: map[string]string{
				"env": "prod",
			},
		},
		{
			name: "multiple_tags",
			input: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"env":  types.StringValue("prod"),
					"team": types.StringValue("core"),
				},
			),
			want: map[string]string{
				"env":  "prod",
				"team": "core",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandTags(ctx, tt.input, &diags)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error, got none")
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			require.Equalf(t, tt.want, got, "expandTags() mismatch")
		})
	}
}

func TestMergeTags(t *testing.T) {
	tests := []struct {
		name         string
		defaultTags  map[string]string
		resourceTags map[string]string
		want         map[string]string
	}{
		{
			name: "no_overlap",
			defaultTags: map[string]string{
				"env": "prod",
			},
			resourceTags: map[string]string{
				"app": "api",
			},
			want: map[string]string{
				"env": "prod",
				"app": "api",
			},
		},
		{
			name: "resource_overrides_default",
			defaultTags: map[string]string{
				"env":  "prod",
				"team": "core",
			},
			resourceTags: map[string]string{
				"env": "dev",
			},
			want: map[string]string{
				"env":  "dev",
				"team": "core",
			},
		},
		{
			name:         "only_default_tags",
			defaultTags:  map[string]string{"env": "prod"},
			resourceTags: map[string]string{},
			want:         map[string]string{"env": "prod"},
		},
		{
			name:         "only_resource_tags",
			defaultTags:  map[string]string{},
			resourceTags: map[string]string{"app": "api"},
			want:         map[string]string{"app": "api"},
		},
		{
			name:         "both_empty",
			defaultTags:  map[string]string{},
			resourceTags: map[string]string{},
			want:         map[string]string{},
		},
		{
			name:         "nil_inputs",
			defaultTags:  nil,
			resourceTags: nil,
			want:         map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeTags(tt.defaultTags, tt.resourceTags)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestDiffTags(t *testing.T) {
	tests := []struct {
		name          string
		current       map[string]string
		desired       map[string]string
		wantRemove    []string
		wantAddUpdate map[string]string
	}{
		{
			name:          "no changes",
			current:       map[string]string{"a": "1", "b": "2"},
			desired:       map[string]string{"a": "1", "b": "2"},
			wantRemove:    nil,
			wantAddUpdate: map[string]string{},
		},
		{
			name:          "remove tag",
			current:       map[string]string{"a": "1", "b": "2"},
			desired:       map[string]string{"a": "1"},
			wantRemove:    []string{"b"},
			wantAddUpdate: map[string]string{},
		},
		{
			name:          "add tag",
			current:       map[string]string{"a": "1"},
			desired:       map[string]string{"a": "1", "b": "2"},
			wantRemove:    nil,
			wantAddUpdate: map[string]string{"b": "2"},
		},
		{
			name:          "update tag value",
			current:       map[string]string{"a": "1"},
			desired:       map[string]string{"a": "2"},
			wantRemove:    nil,
			wantAddUpdate: map[string]string{"a": "2"},
		},
		{
			name:          "add and remove",
			current:       map[string]string{"a": "1", "b": "2"},
			desired:       map[string]string{"a": "1", "c": "3"},
			wantRemove:    []string{"b"},
			wantAddUpdate: map[string]string{"c": "3"},
		},
		{
			name:          "update and add",
			current:       map[string]string{"a": "1"},
			desired:       map[string]string{"a": "2", "b": "3"},
			wantRemove:    nil,
			wantAddUpdate: map[string]string{"a": "2", "b": "3"},
		},
		{
			name:          "empty current",
			current:       map[string]string{},
			desired:       map[string]string{"a": "1"},
			wantRemove:    nil,
			wantAddUpdate: map[string]string{"a": "1"},
		},
		{
			name:          "empty desired",
			current:       map[string]string{"a": "1"},
			desired:       map[string]string{},
			wantRemove:    []string{"a"},
			wantAddUpdate: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remove, addUpdate := diffTags(tt.current, tt.desired)

			assert.ElementsMatch(t, tt.wantRemove, remove, "unexpected tags to remove")

			assert.Equal(t, tt.wantAddUpdate, addUpdate, "unexpected tags to add/update")
		})
	}
}
