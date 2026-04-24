// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
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

// TagManager reconciles desired Terraform tags with AWS tags for a resource.
type TagManager struct {
	client      taggingAPI
	defaultTags map[string]string
}

// NewTagManager creates a TagManager with the given client and default tags.
func NewTagManager(taggingAPI taggingAPI, defaultTags map[string]string) *TagManager {
	return &TagManager{taggingAPI, defaultTags}
}

// ReadAll reads all resource tags for the given ARN.
// Optional AWS SDK options can be provided (for example a region override).
func (tm *TagManager) ReadAll(ctx context.Context, arn string, optFns ...func(*awstaggingapi.Options),
) (types.Map, diag.Diagnostics) {
	tags, diags := tm.readRaw(ctx, arn, optFns...)
	if diags.HasError() {
		return types.MapNull(types.StringType), diags
	}

	return flattenTags(ctx, tags, &diags), diags
}

func (tm *TagManager) readRaw(
	ctx context.Context,
	arn string,
	optFns ...func(*awstaggingapi.Options),
) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if arn == "" {
		return nil, diags
	}

	raw := make(map[string]string)

	out, err := tm.client.GetResources(
		ctx,
		&awstaggingapi.GetResourcesInput{
			ResourceARNList: []string{arn},
		},
		optFns...,
	)
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

// Apply reconciles the desired tags against the current remote tags.
// Optional AWS SDK options can be provided (for example a region override).
func (tm *TagManager) Apply(
	ctx context.Context,
	arn string,
	desired types.Map,
	optFns ...func(*awstaggingapi.Options),
) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	if arn == "" {
		return types.MapNull(types.StringType), diags
	}

	if desired.IsUnknown() {
		// preserve current remote state
		return tm.ReadAll(ctx, arn, optFns...)
	}

	current, readDiags := tm.readRaw(ctx, arn, optFns...)
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
		err := util.RetryOn(
			ctx,
			func(ctx context.Context) error {
				out, err := tm.client.UntagResources(
					ctx,
					&awstaggingapi.UntagResourcesInput{
						ResourceARNList: []string{arn},
						TagKeys:         removeKeys,
					},
					optFns...,
				)
				if err != nil {
					return err
				}

				if out != nil && len(out.FailedResourcesMap) > 0 {
					for _, failedResource := range out.FailedResourcesMap {
						return &smithy.GenericAPIError{
							Code:    string(failedResource.ErrorCode),
							Message: aws.ToString(failedResource.ErrorMessage),
						}
					}
				}

				return err
			},
			util.WithTimeout(taggingRetryTimeout),
			util.WithInitBackoff(taggingRetryInitBackoff),
			util.WithMaxBackoff(taggingRetryMaxBackoff),
			// see https://docs.aws.amazon.com/resourcegroupstagging/latest/APIReference/API_UntagResources.html
			util.WithRetryOnFns(
				util.IsThrottledException,
			),
		)

		if err != nil {
			diags.AddError(
				"Error Removing AWS Tags",
				fmt.Sprintf("Could not remove tags from resource %q: %v", arn, err),
			)
			return types.MapNull(types.StringType), diags
		}
	}

	if len(addOrUpdate) > 0 {
		err := util.RetryOn(
			ctx,
			func(ctx context.Context) error {
				out, err := tm.client.TagResources(
					ctx,
					&awstaggingapi.TagResourcesInput{
						ResourceARNList: []string{arn},
						Tags:            addOrUpdate,
					},
					optFns...,
				)
				if err != nil {
					return err
				}

				if out != nil && len(out.FailedResourcesMap) > 0 {
					for _, failedResource := range out.FailedResourcesMap {
						return &smithy.GenericAPIError{
							Code:    string(failedResource.ErrorCode),
							Message: aws.ToString(failedResource.ErrorMessage),
						}
					}
				}

				return err
			},
			util.WithTimeout(taggingRetryTimeout),
			util.WithInitBackoff(taggingRetryInitBackoff),
			util.WithMaxBackoff(taggingRetryMaxBackoff),
			// see https://docs.aws.amazon.com/resourcegroupstagging/latest/APIReference/API_TagResources.html
			util.WithRetryOnFns(
				util.IsThrottledException,
			),
		)

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

func diffTags(current, desired map[string]string) (removeKeys []string, addOrUpdate map[string]string) {
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

// ResourceTags derives user-managed tags from tags_all and prior state.
func ResourceTags(ctx context.Context, prior types.Map, all types.Map) (resourceTags types.Map, diags diag.Diagnostics) {
	// no tags exist, nothing to return
	if all.IsNull() || all.IsUnknown() {
		resourceTags = types.MapNull(types.StringType)
		return
	}

	allExpanded := expandTags(ctx, all, &diags)
	if diags.HasError() {
		return
	}

	// user never set tags, keep tags null
	if prior.IsNull() || prior.IsUnknown() {
		resourceTags = types.MapNull(types.StringType)
		return
	}

	priorExpanded := expandTags(ctx, prior, &diags)
	if diags.HasError() {
		return
	}

	// user explicitly set empty tags
	if len(priorExpanded) == 0 {
		resourceTags = types.MapValueMust(types.StringType, map[string]attr.Value{})
		return
	}

	filtered := make(map[string]string)
	for k := range priorExpanded {
		if v, ok := allExpanded[k]; ok {
			filtered[k] = v
		}
	}

	resourceTags = flattenTags(ctx, filtered, &diags)
	return
}

// EffectiveTagsForPlan returns the planned effective tags (resource tags merged with provider default tags).
// If tags are unknown, or any tag value is unknown, the result is unknown.
func (tm *TagManager) EffectiveTagsForPlan(resourceTags types.Map) types.Map {
	if resourceTags.IsUnknown() {
		return types.MapUnknown(types.StringType)
	}

	merged := make(map[string]attr.Value)
	for k, v := range tm.defaultTags {
		merged[k] = types.StringValue(v)
	}

	if resourceTags.IsNull() {
		if len(merged) == 0 {
			return types.MapNull(types.StringType)
		}
		return types.MapValueMust(types.StringType, merged)
	}

	for k, v := range resourceTags.Elements() {
		tagValue, ok := v.(types.String)
		if !ok || tagValue.IsUnknown() || tagValue.IsNull() {
			return types.MapUnknown(types.StringType)
		}
		merged[k] = types.StringValue(tagValue.ValueString())
	}

	if len(merged) == 0 {
		return types.MapNull(types.StringType)
	}

	return types.MapValueMust(types.StringType, merged)
}
