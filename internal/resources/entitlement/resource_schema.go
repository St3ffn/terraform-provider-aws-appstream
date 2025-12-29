// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Entitlement",
		MarkdownDescription: "Manages an AppStream entitlement within a specific AppStream stack. " +
			"Entitlements define which users (based on AWS IAM SAML PrincipalTag attributes) can see and launch applications " +
			"from a stack.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the entitlement.",
				MarkdownDescription: "A synthetic identifier for the entitlement, composed of the stack name and entitlement name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"stack_name": schema.StringAttribute{
				Description: "Name of the AppStream Stack.",
				MarkdownDescription: "The name of the AppStream stack in which the entitlement is defined. " +
					"Changing this value forces the entitlement to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream Entitlement.",
				MarkdownDescription: "The name of the entitlement within the stack. " +
					"Changing this value forces the entitlement to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the entitlement.",
				MarkdownDescription: "An optional description for the entitlement. " +
					"Must be 256 characters or fewer.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"app_visibility": schema.StringAttribute{
				Description: "Visibility of applications for this entitlement.",
				MarkdownDescription: "Controls which applications are visible to users who match the entitlement attributes. " +
					"Valid values are `ALL` or `ASSOCIATED`.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ALL", "ASSOCIATED"),
				},
			},
			"attributes": schema.SetNestedAttribute{
				Description: "Entitlement attribute used to match federated user sessions.",
				MarkdownDescription: "Exactly one attribute rule used to match federated user attributes " +
					"(AWS IAM SAML PrincipalTag). Each entitlement supports **one** attribute only. " +
					"To apply multiple rules, create multiple entitlements.",
				Required: true,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Attribute name.",
							MarkdownDescription: "A supported AWS IAM SAML PrincipalTag attribute name. " +
								"Valid values are: `roles`, `department`, `organization`, `groups`, `title`, `costCenter`, `userType`.",
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.OneOf(
									"roles",
									"department",
									"organization",
									"groups",
									"title",
									"costCenter",
									"userType",
								),
							},
						},
						"value": schema.StringAttribute{
							Description: "Attribute value.",
							MarkdownDescription: "The value of the selected attribute name that must match the federated user session. " +
								"Must be at least 1 character.",
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"created_time": schema.StringAttribute{
				Description: "Time the entitlement was created.",
				MarkdownDescription: "The timestamp when the entitlement was created, in RFC 3339 format (for example, `2024-01-09T14:32:11Z`). " +
					"This value is set by AWS and cannot be modified.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
