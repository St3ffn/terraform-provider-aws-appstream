// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Directory Config",
		MarkdownDescription: "Manages an AppStream directory configuration. " +
			"A directory config defines how AppStream fleets and image builders " +
			"join a Microsoft Active Directory domain, including organizational units " +
			"and authentication settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream directory config.",
				MarkdownDescription: "A synthetic identifier for the directory config, equal to the directory name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"directory_name": schema.StringAttribute{
				Description: "Directory name.",
				MarkdownDescription: "The fully qualified domain name of the Microsoft Active Directory " +
					"(for example, `corp.example.com`). Changing this value forces the directory config to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"organizational_unit_distinguished_names": schema.SetAttribute{
				Description: "Organizational unit distinguished names.",
				MarkdownDescription: "The distinguished names of the organizational units used for computer accounts. " +
					"These OUs must allow computer objects to be created and managed by the AppStream service account.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(2000),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
					),
				},
			},
			"service_account_credentials": schema.SingleNestedAttribute{
				Description: "Service account credentials.",
				MarkdownDescription: "Specifies the credentials of the Active Directory service account used " +
					"by AppStream fleets and image builders to join the domain. " +
					"These credentials are write-only and are not returned by AWS after creation.",
				Optional:  true,
				Sensitive: true,
				Attributes: map[string]schema.Attribute{
					"account_name": schema.StringAttribute{
						Description: "Service account user name.",
						MarkdownDescription: "The user name of the Active Directory service account. " +
							"This account must have permissions to create computer objects, join computers to the domain, " +
							"and reset passwords for computer objects in the specified organizational units.",
						Required:  true,
						Sensitive: true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"account_password": schema.StringAttribute{
						Description: "Service account password.",
						MarkdownDescription: "The password for the Active Directory service account. " +
							"This value is sensitive and is never returned by AWS.",
						Required:  true,
						Sensitive: true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 127),
						},
					},
				},
			},
			"certificate_based_auth_properties": schema.SingleNestedAttribute{
				Description: "Certificate-based authentication configuration.",
				MarkdownDescription: "Specifies certificate-based authentication settings used to authenticate " +
					"SAML 2.0 identity provider users to Active Directory domain-joined streaming instances.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Description: "Certificate-based authentication status.",
						MarkdownDescription: "Controls whether certificate-based authentication is enabled. " +
							"Valid values are `DISABLED`, `ENABLED`, or `ENABLED_NO_DIRECTORY_LOGIN_FALLBACK`.",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"DISABLED",
								"ENABLED",
								"ENABLED_NO_DIRECTORY_LOGIN_FALLBACK",
							),
						},
					},
					"certificate_authority_arn": schema.StringAttribute{
						Description: "Private CA ARN.",
						MarkdownDescription: "The ARN of the AWS Certificate Manager Private Certificate Authority " +
							"used for certificate-based authentication.",
						Optional: true,
						Validators: []validator.String{
							util.ValidARNWithServiceAndResource("acm-pca", "certificate-authority/"),
						},
					},
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the directory config was created.",
				MarkdownDescription: "The timestamp when the directory configuration was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
