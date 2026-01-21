// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package sessions

import "github.com/hashicorp/terraform-plugin-framework/types"

type model struct {
	// StackName is the name of the AppStream stack (required).
	StackName types.String `tfsdk:"stack_name"`
	// FleetName is the name of the AppStream fleet (required).
	FleetName types.String `tfsdk:"fleet_name"`
	// UserID is the identifier of the user for whom sessions are listed (optional).
	// If specified, AuthenticationType must also be provided.
	UserID types.String `tfsdk:"user_id"`
	// AuthenticationType is the authentication method used by the user (optional).
	AuthenticationType types.String `tfsdk:"authentication_type"`
	// InstanceID is the identifier of the instance hosting the session (optional).
	InstanceID types.String `tfsdk:"instance_id"`
	// Sessions is the list of streaming sessions that match the specified filters (computed).
	Sessions types.Set `tfsdk:"sessions"`
}

type sessionModel struct {
	// ID is the unique identifier of the streaming session (computed).
	ID types.String `tfsdk:"id"`
	// UserID is the identifier of the user for whom the session was created (computed).
	UserID types.String `tfsdk:"user_id"`
	// StackName is the name of the AppStream stack for the session (computed).
	StackName types.String `tfsdk:"stack_name"`
	// FleetName is the name of the AppStream fleet for the session (computed).
	FleetName types.String `tfsdk:"fleet_name"`
	// State is the current state of the streaming session (computed).
	State types.String `tfsdk:"state"`
	// ConnectionState indicates whether the user is connected to the session (computed).
	ConnectionState types.String `tfsdk:"connection_state"`
	// StartTime is the timestamp when the session started (computed).
	StartTime types.String `tfsdk:"start_time"`
	// MaxExpirationTime is the timestamp when the session is scheduled to expire (computed).
	MaxExpirationTime types.String `tfsdk:"max_expiration_time"`
	// AuthenticationType is the authentication method used for the session (computed).
	AuthenticationType types.String `tfsdk:"authentication_type"`
	// InstanceID is the identifier of the instance hosting the session (computed).
	InstanceID types.String `tfsdk:"instance_id"`
	// NetworkAccessConfiguration contains network details for the session (computed).
	NetworkAccessConfiguration types.Object `tfsdk:"network_access_configuration"`
}

type networkAccessConfigurationModel struct {
	// EniPrivateIPAddress is the private IP address of the elastic network interface (computed).
	EniPrivateIPAddress types.String `tfsdk:"eni_private_ip_address"`
	// EniIPv6Addresses are the IPv6 addresses assigned to the elastic network interface (computed).
	EniIPv6Addresses types.Set `tfsdk:"eni_ipv6_addresses"`
	// EniID is the identifier of the elastic network interface (computed).
	EniID types.String `tfsdk:"eni_id"`
}
