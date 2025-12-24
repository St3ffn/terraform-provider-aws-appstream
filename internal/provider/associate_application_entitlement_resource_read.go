// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *associateApplicationEntitlementResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {

	var state associateApplicationEntitlementModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addAssociateApplicationEntitlementDiagnostics(state, &resp.Diagnostics, associateDiagnosticRead)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readAssociateApplicationEntitlement(
		ctx,
		state.StackName.ValueString(),
		state.EntitlementName.ValueString(),
		state.ApplicationIdentifier.ValueString(),
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if isContextCanceled(ctx.Err()) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *associateApplicationEntitlementResource) readAssociateApplicationEntitlement(
	ctx context.Context, stackName, entitlementName, applicationIdentifier string,
) (*associateApplicationEntitlementModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	var nextToken *string

	for {
		out, err := r.appstreamClient.ListEntitledApplications(ctx, &awsappstream.ListEntitledApplicationsInput{
			StackName:       aws.String(stackName),
			EntitlementName: aws.String(entitlementName),
			NextToken:       nextToken,
			MaxResults:      aws.Int32(AppStreamMaxResults),
		})
		if err != nil {
			if isContextCanceled(err) {
				return nil, diags
			}

			if isAppStreamNotFound(err) {
				return nil, diags
			}

			diags.AddError(
				"Error Reading AWS AppStream Application Entitlement Association",
				fmt.Sprintf(
					"Could not read association of application %q with entitlement %q (stack %q): %v",
					applicationIdentifier, entitlementName, stackName, err,
				),
			)
			return nil, diags
		}

		for _, ea := range out.EntitledApplications {
			if ea.ApplicationIdentifier != nil && *ea.ApplicationIdentifier == applicationIdentifier {
				state := &associateApplicationEntitlementModel{
					ID: types.StringValue(
						buildAssociateApplicationEntitlementID(stackName, entitlementName, applicationIdentifier),
					),
					StackName:             types.StringValue(stackName),
					EntitlementName:       types.StringValue(entitlementName),
					ApplicationIdentifier: types.StringValue(applicationIdentifier),
				}
				return state, diags
			}
		}

		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	return nil, diags
}
