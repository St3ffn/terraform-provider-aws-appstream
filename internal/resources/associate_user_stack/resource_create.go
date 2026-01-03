// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	addDiagnostics(plan, &resp.Diagnostics, diagnosticPlan)
	if resp.Diagnostics.HasError() {
		return
	}

	stackName := plan.StackName.ValueString()
	userName := plan.UserName.ValueString()
	authenticationType := plan.AuthenticationType.ValueString()

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			out, err := r.appstreamClient.BatchAssociateUserStack(ctx, &awsappstream.BatchAssociateUserStackInput{
				UserStackAssociations: []awstypes.UserStackAssociation{
					{
						AuthenticationType:    awstypes.AuthenticationType(authenticationType),
						StackName:             aws.String(stackName),
						UserName:              aws.String(userName),
						SendEmailNotification: util.BoolPointerOrNil(plan.SendEmailNotification),
					},
				},
			})
			if err != nil {
				return err
			}

			var retryErr error

			for _, e := range out.Errors {
				switch e.ErrorCode {
				case awstypes.UserStackAssociationErrorCodeStackNotFound,
					awstypes.UserStackAssociationErrorCodeUserNameNotFound,
					awstypes.UserStackAssociationErrorCodeDirectoryNotFound,
					awstypes.UserStackAssociationErrorCodeInternalError:
					retryErr = newUserStackAssociationNotReadyError(e)
				default:
					// fail on first non-retryable error
					errMsg := "appstream user-stack association failed"
					if e.ErrorMessage != nil {
						errMsg = fmt.Sprintf("%s: %s", errMsg, *e.ErrorMessage)
					}
					return errors.New(errMsg)
				}
			}

			if retryErr != nil {
				return retryErr
			}
			return nil
		},
		util.WithTimeout(3*time.Minute),
		util.WithInitBackoff(1*time.Second),
		util.WithMaxBackoff(30*time.Second),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_BatchAssociateUserStack.html
		util.WithRetryOnFns(
			util.IsOperationNotPermittedException,
			isUserStackAssociationNotReadyError,
		),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream User Stack Association",
			fmt.Sprintf("Could not associate user %q to stack %q: %v",
				userName, stackName, err,
			),
		)
		return
	}

	newState, diags := r.readAssociateUserStack(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if ctx.Err() != nil {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
