// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config model

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if config.StackName.IsNull() || config.StackName.IsUnknown() ||
		config.Name.IsNull() || config.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read entitlement because stack_name and name must be set and known.",
		)
		return
	}

	stackName := config.StackName.ValueString()
	name := config.Name.ValueString()

	out, err := ds.appstreamClient.DescribeEntitlements(ctx, &awsappstream.DescribeEntitlementsInput{
		StackName: aws.String(stackName),
		Name:      aws.String(name),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Entitlement Not Found",
				fmt.Sprintf("No entitlement %q was found in stack %q.", name, stackName),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Entitlement",
			fmt.Sprintf("Could not read entitlement %q in stack %q: %v", name, stackName, err),
		)
		return
	}

	if len(out.Entitlements) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream Entitlement Not Found",
			fmt.Sprintf("No entitlement %q was found in stack %q.", name, stackName),
		)
		return
	}

	e := out.Entitlements[0]
	if e.StackName == nil || e.Name == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("Entitlement %q in stack %q was returned without required identifiers.", name, stackName),
		)
		return
	}

	state := &model{
		ID:            types.StringValue(buildID(aws.ToString(e.StackName), aws.ToString(e.Name))),
		StackName:     types.StringValue(aws.ToString(e.StackName)),
		Name:          types.StringValue(aws.ToString(e.Name)),
		Description:   util.StringOrNull(e.Description),
		AppVisibility: types.StringValue(string(e.AppVisibility)),
		Attributes:    flattenAttributes(ctx, e.Attributes, &resp.Diagnostics),
		CreatedTime:   util.StringFromTime(e.CreatedTime),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
