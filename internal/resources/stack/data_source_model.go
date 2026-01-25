// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// ID is a synthetic identifier composed of "<name>".
	ID types.String `tfsdk:"id"`
	// Name is the name of the AppStream stack (required).
	Name types.String `tfsdk:"name"`
	// Description is the description to display (computed).
	Description types.String `tfsdk:"description"`
	// DisplayName is the stack name to display (computed).
	DisplayName types.String `tfsdk:"display_name"`
	// StorageConnectors is the storage connectors to enable (computed).
	StorageConnectors types.Set `tfsdk:"storage_connectors"`
	// RedirectURL is the URL that users are redirected to after their streaming session ends (computed).
	RedirectURL types.String `tfsdk:"redirect_url"`
	// FeedbackURL is the URL that users are redirected to after they click the Send Feedback link (computed).
	FeedbackURL types.String `tfsdk:"feedback_url"`
	// UserSettings is the actions that are enabled or disabled for users during their streaming sessions (computed).
	UserSettings types.Set `tfsdk:"user_settings"`
	// ApplicationSettings configures application settings persistence for users of this stack (computed).
	ApplicationSettings types.Object `tfsdk:"application_settings"`
	// AccessEndpoints is the list of interface VPC endpoints users of the stack can connect through (computed).
	AccessEndpoints types.Set `tfsdk:"access_endpoints"`
	// EmbedHostDomains is the domains where streaming sessions can be embedded in an iframe (computed).
	EmbedHostDomains types.Set `tfsdk:"embed_host_domains"`
	// StreamingExperienceSettings is the streaming protocol the stack should prefer (computed).
	StreamingExperienceSettings types.Object `tfsdk:"streaming_experience_settings"`
	// Tags is the resource tags to apply to the stack (computed).
	Tags types.Map `tfsdk:"tags"`
	// TagsAll is a map of tags, including default tags, assigned to the stack (computed).
	TagsAll types.Map `tfsdk:"tags_all"`
	// ARN of the AppStream stack (computed).
	ARN types.String `tfsdk:"arn"`
	// CreatedTime is the timestamp when the AppStream stack was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// StackErrors is the list of errors reported by AWS for this stack (computed).
	StackErrors types.Set `tfsdk:"stack_errors"`
}
