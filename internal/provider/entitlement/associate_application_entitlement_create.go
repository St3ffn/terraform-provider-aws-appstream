// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *associateApplicationEntitlementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan associateApplicationEntitlementModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	addAssocPartsDiagnostics(plan, &resp.Diagnostics, assocDiagPlan)
	if resp.Diagnostics.HasError() {
		return
	}

	stackName := plan.StackName.ValueString()
	entName := plan.EntitlementName.ValueString()
	appID := plan.ApplicationIdentifier.ValueString()

	// helper: check if already associated
	isAssociated := func() (bool, error) {
		var nextToken *string
		for {
			out, err := r.appStreamClient.ListEntitledApplications(ctx, &awsappstream.ListEntitledApplicationsInput{
				StackName:       aws.String(stackName),
				EntitlementName: aws.String(entName),
				NextToken:       nextToken,
				MaxResults:      aws.Int32(AppStreamMaxResults),
			})
			if err != nil {
				return false, err
			}
			for _, ea := range out.EntitledApplications {
				if ea.ApplicationIdentifier != nil && *ea.ApplicationIdentifier == appID {
					return true, nil
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				return false, nil
			}
			nextToken = out.NextToken
		}
	}

	already, err := isAssociated()
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}
		if isAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"Error Creating AWS AppStream Application Entitlement Association",
				fmt.Sprintf("Stack or entitlement not found while checking existing association of application %q with entitlement %q (stack %q): %v",
					appID, entName, stackName, err,
				),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Application Entitlement Association",
			fmt.Sprintf("Could not check existing association of application %q with entitlement %q (stack %q): %v",
				appID, entName, stackName, err,
			),
		)
		return
	}

	plan.ID = types.StringValue(buildAssocID(stackName, entName, appID))

	if already {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	_, err = r.appStreamClient.AssociateApplicationToEntitlement(ctx, &awsappstream.AssociateApplicationToEntitlementInput{
		StackName:             aws.String(stackName),
		EntitlementName:       aws.String(entName),
		ApplicationIdentifier: aws.String(appID),
	})
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Application Entitlement Association",
			fmt.Sprintf("Could not associate application %q to entitlement %q (stack %q): %v",
				appID, entName, stackName, err,
			),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func buildAssocID(stack, ent, app string) string {
	return fmt.Sprintf("%s|%s|%s", stack, ent, app)
}
