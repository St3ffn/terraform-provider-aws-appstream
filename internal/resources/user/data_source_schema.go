// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (d *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream User",
		MarkdownDescription: "Reads information about an AppStream user. " +
			"The user is uniquely identified by the user name and authentication type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream user.",
				MarkdownDescription: "A synthetic identifier for the user, composed of the authentication type and user name " +
					"in the format `<authentication_type>|<user_name>`.",
				Computed: true,
			},
			"authentication_type": schema.StringAttribute{
				Description: "Authentication type for the user.",
				MarkdownDescription: "The authentication type associated with the user. " +
					"Valid values are `API`, `SAML`, `USERPOOL`, or `AWS_AD`.",
				Required: true,
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
					"Email addresses are **case-sensitive**.",
				Required: true,
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
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				Description:         "Last name of the user.",
				MarkdownDescription: "The last (family) name of the user.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "Whether the user is enabled.",
				MarkdownDescription: "Indicates whether the user is enabled.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				Description:         "User status.",
				MarkdownDescription: "The status of the user as reported by AWS.",
				Computed:            true,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream user.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream user.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the user was created.",
				MarkdownDescription: "The timestamp when the user was created, in RFC 3339 format.",
				Computed:            true,
			},
		},
	}
}
