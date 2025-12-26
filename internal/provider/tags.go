// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type taggingAPI interface {
	GetResources(
		ctx context.Context, params *awstaggingapi.GetResourcesInput, optFns ...func(*awstaggingapi.Options),
	) (*awstaggingapi.GetResourcesOutput, error)

	TagResources(
		ctx context.Context, params *awstaggingapi.TagResourcesInput, optFns ...func(*awstaggingapi.Options),
	) (*awstaggingapi.TagResourcesOutput, error)

	UntagResources(
		ctx context.Context, params *awstaggingapi.UntagResourcesInput, optFns ...func(*awstaggingapi.Options),
	) (*awstaggingapi.UntagResourcesOutput, error)
}

type tagManager struct {
	client      taggingAPI
	defaultTags map[string]string
}

func newTagManager(taggingAPI taggingAPI, defaultTags map[string]string) *tagManager {
	return &tagManager{taggingAPI, defaultTags}
}

func (tm *tagManager) Read(ctx context.Context, arn string) (types.Map, diag.Diagnostics) {
	tags, diags := tm.readRaw(ctx, arn)
	if diags.HasError() {
		return types.MapNull(types.StringType), diags
	}

	return flattenTags(ctx, tags, &diags), diags
}

func (tm *tagManager) readRaw(ctx context.Context, arn string) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if arn == "" {
		return nil, diags
	}

	raw := make(map[string]string)

	out, err := tm.client.GetResources(ctx, &awstaggingapi.GetResourcesInput{
		ResourceARNList: []string{arn},
	})
	if err != nil {
		diags.AddError(
			"Error Reading AWS Tags",
			fmt.Sprintf("Could not read tags for resource %q: %v", arn, err),
		)
		return nil, diags
	}

	for _, m := range out.ResourceTagMappingList {
		for _, t := range m.Tags {
			if t.Key != nil && t.Value != nil {
				raw[*t.Key] = *t.Value
			}
		}
	}

	return raw, diags
}

func (tm *tagManager) Apply(ctx context.Context, arn string, desired types.Map) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	if arn == "" {
		return types.MapNull(types.StringType), diags
	}

	if desired.IsUnknown() {
		// preserve current remote state
		return tm.Read(ctx, arn)
	}

	current, readDiags := tm.readRaw(ctx, arn)
	diags.Append(readDiags...)
	if diags.HasError() {
		return types.MapNull(types.StringType), diags
	}

	desiredTags := tm.defaultTags
	if !desired.IsNull() {
		resourceTags := expandTags(ctx, desired, &diags)
		if diags.HasError() {
			return types.MapNull(types.StringType), diags
		}
		desiredTags = mergeTags(tm.defaultTags, resourceTags)
	}

	removeKeys, addOrUpdate := diffTags(current, desiredTags)

	if len(removeKeys) > 0 {
		_, err := tm.client.UntagResources(ctx, &awstaggingapi.UntagResourcesInput{
			ResourceARNList: []string{arn},
			TagKeys:         removeKeys,
		})
		if err != nil {
			diags.AddError(
				"Error Removing AWS Tags",
				fmt.Sprintf("Could not remove tags from resource %q: %v", arn, err),
			)
			return types.MapNull(types.StringType), diags
		}
	}

	if len(addOrUpdate) > 0 {
		_, err := tm.client.TagResources(ctx, &awstaggingapi.TagResourcesInput{
			ResourceARNList: []string{arn},
			Tags:            addOrUpdate,
		})
		if err != nil {
			diags.AddError(
				"Error Updating AWS Tags",
				fmt.Sprintf("Could not update tags for resource %q: %v", arn, err),
			)
			return types.MapNull(types.StringType), diags
		}
	}

	return flattenTags(ctx, desiredTags, &diags), diags
}

func flattenTags(ctx context.Context, tags map[string]string, diags *diag.Diagnostics) types.Map {
	if len(tags) == 0 {
		return types.MapNull(types.StringType)
	}

	m, d := types.MapValueFrom(ctx, types.StringType, tags)
	diags.Append(d...)
	if diags.HasError() {
		return types.MapNull(types.StringType)
	}

	return m
}

func expandTags(ctx context.Context, m types.Map, diags *diag.Diagnostics) map[string]string {
	var tags map[string]string
	diags.Append(m.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil
	}

	return tags
}

func mergeTags(defaultTags, resourceTags map[string]string) map[string]string {
	out := make(map[string]string)

	for k, v := range defaultTags {
		out[k] = v
	}
	for k, v := range resourceTags {
		out[k] = v
	}
	return out
}

func diffTags(
	current map[string]string, desired map[string]string,
) (removeKeys []string, addOrUpdate map[string]string) {

	addOrUpdate = make(map[string]string)

	for k, v := range current {
		if desiredVal, ok := desired[k]; !ok {
			removeKeys = append(removeKeys, k)
		} else if desiredVal != v {
			// tag value changed
			addOrUpdate[k] = desiredVal
		}
	}

	for k, v := range desired {
		if _, ok := current[k]; !ok {
			addOrUpdate[k] = v
		}
	}

	return removeKeys, addOrUpdate
}
