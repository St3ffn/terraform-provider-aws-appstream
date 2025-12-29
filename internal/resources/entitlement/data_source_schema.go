// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream Entitlement",
		MarkdownDescription: "Reads an AppStream entitlement within a specific AppStream stack. " +
			"An entitlement defines a single attribute rule that controls which users can see and launch applications " +
			"from the stack. This data source can be used to reference entitlements that are managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the entitlement.",
				MarkdownDescription: "A synthetic identifier for the entitlement, composed of the stack name and entitlement name " +
					"in the format `<stack_name>|<name>`.",
				Computed: true,
			},
			"stack_name": schema.StringAttribute{
				Description:         "Name of the AppStream Stack.",
				MarkdownDescription: "The name of the AppStream stack in which the entitlement is defined.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the AppStream Entitlement.",
				MarkdownDescription: "The name of the entitlement within the stack.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the entitlement.",
				MarkdownDescription: "The entitlement description, if set.",
				Computed:            true,
			},
			"app_visibility": schema.StringAttribute{
				Description: "Visibility of applications for this entitlement.",
				MarkdownDescription: "Controls which applications are visible to users who match the entitlement attributes. " +
					"Valid values are `ALL` or `ASSOCIATED`.",
				Computed: true,
			},
			"attributes": schema.SetNestedAttribute{
				Description: "Entitlement attribute used to match federated user sessions.",
				MarkdownDescription: "The attribute rule used to match federated user attributes (AWS IAM SAML PrincipalTag). " +
					"Each entitlement supports **exactly one** attribute. " +
					"To apply multiple rules, multiple entitlements must be created.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Attribute name.",
							MarkdownDescription: "A supported AWS IAM SAML PrincipalTag attribute name. " +
								"Valid values are: `roles`, `department`, `organization`, `groups`, `title`, `costCenter`, `userType`.",
							Computed: true,
						},
						"value": schema.StringAttribute{
							Description:         "Attribute value.",
							MarkdownDescription: "The value of the selected attribute name that must match the federated user session.",
							Computed:            true,
						},
					},
				},
			},
			"created_time": schema.StringAttribute{
				Description: "Time the entitlement was created.",
				MarkdownDescription: "The timestamp when the entitlement was created, in RFC 3339 format " +
					"(for example, `2024-01-09T14:32:11Z`).",
				Computed: true,
			},
		},
	}
}
