// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *associateApplicationEntitlementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state associateApplicationEntitlementModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addAssocPartsDiagnostics(state, &resp.Diagnostics, assocDiagRead)
	if resp.Diagnostics.HasError() {
		return
	}

	stackName := state.StackName.ValueString()
	entName := state.EntitlementName.ValueString()
	appID := state.ApplicationIdentifier.ValueString()

	var nextToken *string
	found := false

	for {
		out, err := r.appStreamClient.ListEntitledApplications(ctx, &awsappstream.ListEntitledApplicationsInput{
			StackName:       aws.String(stackName),
			EntitlementName: aws.String(entName),
			NextToken:       nextToken,
			MaxResults:      aws.Int32(AppStreamMaxResults),
		})
		if err != nil {
			// respect cancellation/deadlines
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			if isAppStreamNotFound(err) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(
				"Error Reading AWS AppStream Application Entitlement Association",
				fmt.Sprintf(
					"Could not verify association of application %q with entitlement %q (stack %q): %v",
					appID, entName, stackName, err,
				),
			)
			return
		}

		for _, ea := range out.EntitledApplications {
			if ea.ApplicationIdentifier != nil && *ea.ApplicationIdentifier == appID {
				found = true
				break
			}
		}

		if found || out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	if !found {
		// remove resource if missing
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
