// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import "github.com/hashicorp/terraform-plugin-framework/types"

type model struct {
	// ID is a synthetic identifier composed of "<directory_name>" (computed).
	ID types.String `tfsdk:"id"`
	// DirectoryName is the fully qualified domain name of the Microsoft Active Directory (required).
	DirectoryName types.String `tfsdk:"directory_name"`
	// OrganizationalUnitDistinguishedNames are the distinguished names of the organizational units
	// used for computer accounts (required).
	OrganizationalUnitDistinguishedNames types.Set `tfsdk:"organizational_unit_distinguished_names"`
	// ServiceAccountCredentials specifies the credentials of the Active Directory service account
	// used by AppStream fleets and image builders to join the domain (optional, sensitive).
	ServiceAccountCredentials types.Object `tfsdk:"service_account_credentials"`
	// CertificateBasedAuthProperties specifies certificate-based authentication settings
	// for domain-joined streaming instances (optional).
	CertificateBasedAuthProperties types.Object `tfsdk:"certificate_based_auth_properties"`
	// CreatedTime is the timestamp when the directory configuration was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
}

type serviceAccountCredentialsModel struct {
	// AccountName is the user name of the Active Directory service account (required).
	AccountName types.String `tfsdk:"account_name"`
	// AccountPassword is the password of the Active Directory service account (required, sensitive).
	AccountPassword types.String `tfsdk:"account_password"`
}

type certificateBasedAuthPropertiesModel struct {
	// Status controls whether certificate-based authentication is enabled (optional).
	Status types.String `tfsdk:"status"`
	// CertificateAuthorityARN is the ARN of the ACM Private Certificate Authority
	// used for certificate-based authentication (optional).
	CertificateAuthorityARN types.String `tfsdk:"certificate_authority_arn"`
}
