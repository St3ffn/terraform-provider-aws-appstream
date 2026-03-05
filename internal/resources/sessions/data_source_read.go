// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package sessions

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Read fetches the current remote object and refreshes Terraform state from AWS.
// When the object no longer exists remotely, the resource is removed from state
// to converge Terraform with external deletions.
func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if config.StackName.IsNull() || config.StackName.IsUnknown() ||
		config.FleetName.IsNull() || config.FleetName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read sessions because stack_name and fleet_name must be set and known.",
		)
		return
	}

	stackName := config.StackName.ValueString()
	fleetName := config.FleetName.ValueString()

	var sessions []awstypes.Session
	var nextToken *string

	for {
		input := &awsappstream.DescribeSessionsInput{
			FleetName:  aws.String(fleetName),
			StackName:  aws.String(stackName),
			UserId:     util.StringPointerOrNil(config.UserID),
			InstanceId: util.StringPointerOrNil(config.InstanceID),
			NextToken:  nextToken,
		}

		if !config.AuthenticationType.IsNull() && !config.AuthenticationType.IsUnknown() {
			input.AuthenticationType = awstypes.AuthenticationType(config.AuthenticationType.ValueString())
		}

		out, err := ds.appstreamClient.DescribeSessions(ctx, input)
		if err != nil {
			if util.IsContextCanceled(err) {
				return
			}

			resp.Diagnostics.AddError(
				"Error Reading AWS AppStream Sessions",
				fmt.Sprintf("Could not read sessions for stack %q and fleet %q: %v", stackName, fleetName, err),
			)
			return
		}

		sessions = append(sessions, out.Sessions...)

		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	state := &dataSourceModel{
		StackName:          types.StringValue(stackName),
		FleetName:          types.StringValue(fleetName),
		UserID:             config.UserID,
		AuthenticationType: config.AuthenticationType,
		InstanceID:         config.InstanceID,
		Sessions:           flattenSessions(ctx, sessions, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
