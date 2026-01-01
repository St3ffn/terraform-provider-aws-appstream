// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream Directory Config",
		MarkdownDescription: "Reads an AppStream directory config. " +
			"This data source can be used to reference an existing AppStream directory configuration " +
			"that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream directory config.",
				MarkdownDescription: "A synthetic identifier for the directory config, equal to the directory config name.",
				Computed:            true,
			},
			"directory_name": schema.StringAttribute{
				Description: "Directory name.",
				MarkdownDescription: "The fully qualified domain name of the Microsoft Active Directory " +
					"(for example, `corp.example.com`).",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"organizational_unit_distinguished_names": schema.SetAttribute{
				Description: "Organizational unit distinguished names.",
				MarkdownDescription: "The distinguished names of the organizational units " +
					"used for computer accounts in the directory.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"service_account_credentials": schema.SingleNestedAttribute{
				Description: "Service account credentials.",
				MarkdownDescription: "The credentials for the service account used by AppStream " +
					"fleets or image builders to join resources to the Active Directory domain. " +
					"Credential values are sensitive and may not be returned by AWS.",
				Computed:  true,
				Sensitive: true,
				Attributes: map[string]schema.Attribute{
					"account_name": schema.StringAttribute{
						Description: "Service account name.",
						MarkdownDescription: "The user name of the service account used to join resources to the directory. " +
							"This value is sensitive and may not be returned by AWS.",
						Computed:  true,
						Sensitive: true,
					},
					"account_password": schema.StringAttribute{
						Description: "Service account password.",
						MarkdownDescription: "The password of the service account used to join resources to the directory. " +
							"This value is sensitive and is not returned by AWS.",
						Computed:  true,
						Sensitive: true,
					},
				},
			},
			"certificate_based_auth_properties": schema.SingleNestedAttribute{
				Description: "Certificate-based authentication properties.",
				MarkdownDescription: "Configuration for certificate-based authentication " +
					"used to authenticate SAML 2.0 identity provider users to Active Directory " +
					"domain-joined streaming instances, if configured.",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Description: "Certificate-based authentication status.",
						MarkdownDescription: "The status of certificate-based authentication. " +
							"Valid values are `DISABLED`, `ENABLED`, or `ENABLED_NO_DIRECTORY_LOGIN_FALLBACK`.",
						Computed: true,
					},
					"certificate_authority_arn": schema.StringAttribute{
						Description: "Certificate authority ARN.",
						MarkdownDescription: "The ARN of the AWS Certificate Manager Private Certificate Authority " +
							"used for certificate-based authentication, if configured.",
						Computed: true,
					},
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the directory config was created.",
				MarkdownDescription: "The timestamp when the directory configuration was created, in RFC 3339 format.",
				Computed:            true,
			},
		},
	}
}
