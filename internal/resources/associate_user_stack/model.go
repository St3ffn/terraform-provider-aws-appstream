// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type model struct {
	// ID is a synthetic identifier composed of
	// "<stack_name>|<authentication_type>|<user_name>".
	ID types.String `tfsdk:"id"`
	// StackName is the name of the AppStream stack to associate with the user (required).
	StackName types.String `tfsdk:"stack_name"`
	// UserName is the email address of the AppStream user (required).
	// Email addresses are case-sensitive.
	UserName types.String `tfsdk:"user_name"`
	// AuthenticationType is the authentication type for the user (required).
	// Valid values are API, SAML, USERPOOL, or AWS_AD.
	AuthenticationType types.String `tfsdk:"authentication_type"`
	// SendEmailNotification specifies whether a welcome email is sent to the user.
	// This attribute is only used during creation and is not persisted by AWS.
	SendEmailNotification types.Bool `tfsdk:"send_email_notification"`
}
