// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type stackModel struct {
	// ID is a synthetic identifier composed of "<name>".
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream stack (required).
	Name types.String `tfsdk:"name"`
	// Description is the description to display (optional).
	Description types.String `tfsdk:"description"`
	// DisplayName is the stack name to display (optional).
	DisplayName types.String `tfsdk:"display_name"`
	// StorageConnectors is the storage connectors to enable (optional).
	StorageConnectors types.Set `tfsdk:"storage_connectors"`
	// RedirectURL is the URL that users are redirected to after their streaming session ends (optional).
	RedirectURL types.String `tfsdk:"redirect_url"`
	// FeedbackURL is the URL that users are redirected to after they click the Send Feedback link (optional).
	FeedbackURL types.String `tfsdk:"feedback_url"`
	// UserSettings is the actions that are enabled or disabled for users during their streaming sessions (optional).
	UserSettings types.Set `tfsdk:"user_settings"`
	// ApplicationSettings configures application settings persistence for users of this stack (optional).
	ApplicationSettings types.Object `tfsdk:"application_settings"`
	// Tags is the resource tags to apply to the stack (optional).
	Tags types.Map `tfsdk:"tags"`
	// AccessEndpoints is the list of interface VPC endpoints users of the stack can connect through (optional).
	AccessEndpoints types.Set `tfsdk:"access_endpoints"`
	// EmbedHostDomains is the domains where streaming sessions can be embedded in an iframe (optional).
	EmbedHostDomains types.Set `tfsdk:"embed_host_domains"`
	// StreamingExperienceSettings is the streaming protocol the stack should prefer (optional).
	StreamingExperienceSettings types.Object `tfsdk:"streaming_experience_settings"`
	// ARN of the AppStream stack (computed).
	ARN types.String `tfsdk:"arn"`
	// CreatedTime is the timestamp when the AppStream stack was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// StackErrors is the list of errors reported by AWS for this stack (computed).
	StackErrors types.Set `tfsdk:"stack_errors"`
}

type stackStorageConnectorModel struct {
	// ConnectorType is the type of storage connector (required).
	ConnectorType types.String `tfsdk:"connector_type"`
	// ResourceIdentifier is the ARN of the storage connector (optional).
	ResourceIdentifier types.String `tfsdk:"resource_identifier"`
	// Domains is the names of the domains for the account (optional).
	Domains types.Set `tfsdk:"domains"`
	// DomainsRequireAdminConsent is the OneDrive domains where admin consent is required (optional).
	DomainsRequireAdminConsent types.Set `tfsdk:"domains_require_admin_consent"`
}

type stackUserSettingModel struct {
	// Action is the action that can be enabled or disabled for users (required).
	Action types.String `tfsdk:"action"`
	// Permission indicates whether the action is enabled or disabled (required).
	Permission types.String `tfsdk:"permission"`
	// MaximumLength specifies the maximum number of characters that can be copied
	// for clipboard-related actions (optional).
	MaximumLength types.Int64 `tfsdk:"maximum_length"`
}

type stackApplicationSettingsModel struct {
	// Enabled enables application settings persistence (required).
	Enabled types.Bool `tfsdk:"enabled"`
	// SettingsGroup is the name of the settings group (optional).
	SettingsGroup types.String `tfsdk:"settings_group"`
	// S3BucketName is the S3 bucket used for persistent application settings (computed).
	S3BucketName types.String `tfsdk:"s3_bucket_name"`
}

type stackAccessEndpointModel struct {
	// EndpointType is the type of interface endpoint (required).
	EndpointType types.String `tfsdk:"endpoint_type"`
	// VpceID is the ID of the interface VPC endpoint (optional).
	VpceID types.String `tfsdk:"vpce_id"`
}

type stackStreamingExperienceSettingsModel struct {
	// PreferredProtocol is the preferred streaming protocol (optional).
	PreferredProtocol types.String `tfsdk:"preferred_protocol"`
}

type stackErrorModel struct {
	// ErrorCode is the error code reported by AWS for the stack (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message reported by AWS (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
}
