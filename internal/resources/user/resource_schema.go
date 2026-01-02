// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream User",
		MarkdownDescription: "Manages an AppStream user. " +
			"An AppStream user represents an identity that can authenticate to AppStream using a specific authentication type. " +
			"Users can be enabled or disabled and are associated with exactly one authentication type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream user.",
				MarkdownDescription: "A synthetic identifier for the user, composed of the authentication type and user name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authentication_type": schema.StringAttribute{
				Description: "Authentication type for the user.",
				MarkdownDescription: "The authentication type used by the user. " +
					"Changing this value forces the user to be replaced. " +
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
			"user_name": schema.StringAttribute{
				Description: "User name (email address).",
				MarkdownDescription: "The email address of the AppStream user. " +
					"Email addresses are **case-sensitive** and must exactly match during login. " +
					"Changing this value forces the user to be replaced.",
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
			"first_name": schema.StringAttribute{
				Description:         "First name of the user.",
				MarkdownDescription: "The first (given) name of the user.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(2048),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9_\-\s]+$`),
						"must match ^[A-Za-z0-9_\\-\\s]+$",
					),
				},
			},
			"last_name": schema.StringAttribute{
				Description:         "Last name of the user.",
				MarkdownDescription: "The last (family) name of the user.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(2048),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9_\-\s]+$`),
						"must match ^[A-Za-z0-9_\\-\\s]+$",
					),
				},
			},
			"message_action": schema.StringAttribute{
				Description: "Welcome email behavior.",
				MarkdownDescription: "Controls the welcome email sent to the user when the user is created. " +
					"This setting is **write-only** and applies only during creation. " +
					"Valid values are `SUPPRESS` or `RESEND`. " +
					"Changing this value forces the user to be replaced.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"SUPPRESS",
						"RESEND",
					),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the user is enabled.",
				MarkdownDescription: "Specifies whether the user is enabled. " +
					"Defaults to `true` if not explicitly set. " +
					"Disabling a user prevents login without deleting the user.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"status": schema.StringAttribute{
				Description: "User status.",
				MarkdownDescription: "The status of the user as reported by AWS. " +
					"This attribute is informational and cannot be modified.",
				Computed: true,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream user.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream user.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the user was created.",
				MarkdownDescription: "The timestamp when the user was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
