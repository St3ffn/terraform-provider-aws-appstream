// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream User-Stack Association",
		MarkdownDescription: "Manages the association between an AppStream user and an AppStream stack. " +
			"This resource represents the relationship only and does not create or manage the underlying user or stack.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream user-stack association.",
				MarkdownDescription: "A synthetic identifier for the association, composed of the stack name, " +
					"authentication type, and user name. This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"stack_name": schema.StringAttribute{
				Description: "Name of the AppStream stack.",
				MarkdownDescription: "The name of the AppStream stack to associate with the user. " +
					"Changing this value forces the association to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"user_name": schema.StringAttribute{
				Description: "User name (email address).",
				MarkdownDescription: "The email address of the AppStream user to associate with the stack. " +
					"Email addresses are **case-sensitive**. Changing this value forces the association to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`[\p{L}\p{M}\p{S}\p{N}\p{P}]+`),
						"must match [\\p{L}\\p{M}\\p{S}\\p{N}\\p{P}]+",
					),
				},
			},
			"authentication_type": schema.StringAttribute{
				Description: "Authentication type for the user.",
				MarkdownDescription: "The authentication type used by the user. " +
					"Changing this value forces the association to be replaced. " +
					"Valid values are `API`, `SAML`, `USERPOOL`, or `AWS_AD`.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"API",
						"SAML",
						"USERPOOL",
						"AWS_AD",
					),
				},
			},
			"send_email_notification": schema.BoolAttribute{
				Description: "Whether to send a welcome email to the user.",
				MarkdownDescription: "Specifies whether a welcome email is sent to the user after the association is created. " +
					"This option is only applicable when `authentication_type` is `USERPOOL`. " +
					"For other authentication types, the user must already exist and this value is ignored. " +
					"Changing this value forces the association to be replaced.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
