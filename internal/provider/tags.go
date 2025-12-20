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

func readTags(
	ctx context.Context, taggingClient *awstaggingapi.Client, arn string,
) (map[string]string, diag.Diagnostics) {

	var diags diag.Diagnostics
	tags := make(map[string]string)

	if arn == "" {
		return tags, diags
	}

	out, err := taggingClient.GetResources(ctx, &awstaggingapi.GetResourcesInput{
		ResourceARNList: []string{arn},
	})
	if err != nil {
		diags.AddError(
			"Error Reading AWS Tags",
			fmt.Sprintf("Could not read tags for resource %q: %v", arn, err),
		)
		return tags, diags
	}

	for _, m := range out.ResourceTagMappingList {
		for _, t := range m.Tags {
			if t.Key != nil && t.Value != nil {
				tags[*t.Key] = *t.Value
			}
		}
	}

	return tags, diags
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
