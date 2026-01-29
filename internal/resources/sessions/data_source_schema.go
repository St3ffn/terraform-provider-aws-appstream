// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package sessions

import (
	"context"
	"regexp"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read AWS AppStream Sessions",
		MarkdownDescription: "Reads active AppStream streaming sessions for a specific stack and fleet. " +
			"This data source is read-only and can be used for observability, monitoring, and debugging purposes.",
		Attributes: map[string]schema.Attribute{
			"stack_name": schema.StringAttribute{
				Description:         "Name of the AppStream stack.",
				MarkdownDescription: "The name of the AppStream stack for which streaming sessions are listed.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"fleet_name": schema.StringAttribute{
				Description:         "Name of the AppStream fleet.",
				MarkdownDescription: "The name of the AppStream fleet for which streaming sessions are listed.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "User identifier.",
				MarkdownDescription: "The identifier of the user for whom streaming sessions are listed. " +
					"If specified, `authentication_type` must also be provided.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 128),
				},
			},
			"authentication_type": schema.StringAttribute{
				Description: "Authentication type.",
				MarkdownDescription: "The authentication method used by the user. " +
					"Valid values are `API`, `SAML`, `USERPOOL`, or `AWS_AD`.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						util.AWSEnumToSlice(awstypes.AuthenticationType.Values)...,
					),
				},
			},
			"instance_id": schema.StringAttribute{
				Description:         "Instance identifier.",
				MarkdownDescription: "The identifier of the streaming instance hosting the session.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sessions": schema.SetNestedAttribute{
				Description:         "Streaming sessions.",
				MarkdownDescription: "The list of active streaming sessions that match the specified filters.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "Session identifier.",
							MarkdownDescription: "The unique identifier of the streaming session.",
							Computed:            true,
						},
						"user_id": schema.StringAttribute{
							Description:         "User identifier.",
							MarkdownDescription: "The identifier of the user for whom the session was created.",
							Computed:            true,
						},
						"stack_name": schema.StringAttribute{
							Description:         "Stack name.",
							MarkdownDescription: "The name of the AppStream stack associated with the session.",
							Computed:            true,
						},
						"fleet_name": schema.StringAttribute{
							Description:         "Fleet name.",
							MarkdownDescription: "The name of the AppStream fleet associated with the session.",
							Computed:            true,
						},
						"state": schema.StringAttribute{
							Description:         "Session state.",
							MarkdownDescription: "The current state of the streaming session.",
							Computed:            true,
						},
						"connection_state": schema.StringAttribute{
							Description:         "Connection state.",
							MarkdownDescription: "Indicates whether the user is currently connected to the streaming session.",
							Computed:            true,
						},
						"start_time": schema.StringAttribute{
							Description:         "Session start time.",
							MarkdownDescription: "The timestamp when the streaming session started, in RFC 3339 format.",
							Computed:            true,
						},
						"max_expiration_time": schema.StringAttribute{
							Description:         "Session expiration time.",
							MarkdownDescription: "The timestamp when the streaming session is scheduled to expire, in RFC 3339 format.",
							Computed:            true,
						},
						"authentication_type": schema.StringAttribute{
							Description:         "Authentication type.",
							MarkdownDescription: "The authentication method used for the session.",
							Computed:            true,
						},
						"instance_id": schema.StringAttribute{
							Description:         "Instance identifier.",
							MarkdownDescription: "The identifier of the instance hosting the streaming session.",
							Computed:            true,
						},
						"network_access_configuration": schema.SingleNestedAttribute{
							Description:         "Network access configuration.",
							MarkdownDescription: "Network details of the elastic network interface attached to the session.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"eni_private_ip_address": schema.StringAttribute{
									Description:         "Private IP address.",
									MarkdownDescription: "The private IP address of the elastic network interface.",
									Computed:            true,
								},
								"eni_ipv6_addresses": schema.SetAttribute{
									Description:         "IPv6 addresses.",
									MarkdownDescription: "The IPv6 addresses assigned to the elastic network interface.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"eni_id": schema.StringAttribute{
									Description:         "Elastic network interface ID.",
									MarkdownDescription: "The identifier of the elastic network interface.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}
