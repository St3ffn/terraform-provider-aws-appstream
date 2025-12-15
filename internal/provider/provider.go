// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	awsretry "github.com/aws/aws-sdk-go-v2/aws/retry"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	intaws "github.com/st3ffn/terraform-provider-aws-appstream/internal/aws"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/provider/entitlement"
)

var (
	_ provider.Provider = &awsAppStreamProvider{}
	//_ provider.ProviderWithFunctions = &awsAppStreamProvider{} // TODO @MARCUS.
	//_ provider.ProviderWithEphemeralResources = &awsAppStreamProvider{} // TODO @MARCUS.
)

type awsAppStreamProvider struct {
	version string
}

// awsAppStreamProviderModel describes the provider data model.
type awsAppStreamProviderModel struct {
	Profile          types.String `tfsdk:"profile"`
	Region           types.String `tfsdk:"region"`
	RetryMaxAttempts types.Int64  `tfsdk:"retry_max_attempts"`
	RetryMaxBackoff  types.Int64  `tfsdk:"retry_max_backoff"`
}

func (A *awsAppStreamProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.Version = A.version
	resp.TypeName = "awsappstream"
}

func (A *awsAppStreamProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with AWS AppStream",
		MarkdownDescription: `
The **awsappstream** provider allows you to manage AWS AppStream resources
using Terraform.

Authentication and region selection follow the standard AWS SDK behavior.
`,
		Attributes: map[string]schema.Attribute{
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "AWS Profile name. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: "The name of the AWS CLI profile to use. " +
					"If not set, the AWS SDK default credential resolution chain is used (environment variables, shared credentials file, EC2/ECS metadata, etc.).",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"region": schema.StringAttribute{
				Required:            true,
				Description:         "Required AWS Region.",
				MarkdownDescription: "The AWS region in which AppStream resources are managed. This value must be set explicitly and cannot be unknown at plan time.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"retry_max_attempts": schema.Int64Attribute{
				Optional:    true,
				Description: fmt.Sprintf("Maximum number of attempts to retry AWS AppStream operations (default: %d).", awsretry.DefaultMaxAttempts),
				MarkdownDescription: fmt.Sprintf("The maximum number of retry attempts for AWS AppStream API calls. "+
					"This controls how many times a failed request is retried before returning an error. **Default:** %d", awsretry.DefaultMaxAttempts),
				Validators: []validator.Int64{int64validator.AtLeast(1)},
			},
			"retry_max_backoff": schema.Int64Attribute{
				Optional:    true,
				Description: fmt.Sprintf("Maximum retry backoff in seconds for AWS AppStream operations (default: %d).", int64(awsretry.DefaultMaxBackoff.Seconds())),
				MarkdownDescription: fmt.Sprintf("The maximum backoff time, in seconds, between retry attempts for AWS AppStream API calls. "+
					"This value limits the exponential backoff applied by the AWS SDK. **Default:** %d seconds", int64(awsretry.DefaultMaxBackoff.Seconds())),
				Validators: []validator.Int64{int64validator.AtLeast(1)},
			},
		},
	}
}

func (A *awsAppStreamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring AWS AppStream Provider")

	var config awsAppStreamProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Profile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("profile"),
			"Unknown AWS Profile",
			"The AWS AppStream provider cannot be configured because \"profile\" is unknown. "+
				"Provider configuration values must be static. "+
				"Set \"profile\" to a fixed string or remove it to use the AWS SDK default.",
		)
		return
	}

	if config.Region.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Unknown AWS Region",
			"The AWS AppStream provider cannot be configured because \"region\" is unknown. "+
				"\"region\" must be set to a fixed string value.",
		)
		return
	}

	if config.RetryMaxAttempts.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("retry_max_attempts"),
			"Unknown AWS Retry Max Attempts",
			"The AWS AppStream provider cannot be configured because \"retry_max_attempts\" is unknown. "+
				"Provider configuration values must be static. "+
				"Set it to a fixed number or remove it to use the default.",
		)
		return
	}

	if config.RetryMaxBackoff.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("retry_max_backoff"),
			"Unknown AWS Retry Max Backoff",
			"The AWS AppStream provider cannot be configured because \"retry_max_backoff\" is unknown. "+
				"Provider configuration values must be static. "+
				"Set it to a fixed number or remove it to use the default.",
		)
		return
	}

	var awsopts []func(*awsconfig.LoadOptions) error

	if !config.Profile.IsNull() {
		profile := config.Profile.ValueString()
		awsopts = append(awsopts, awsconfig.WithSharedConfigProfile(profile))
		ctx = tflog.SetField(ctx, "profile", profile)
	}

	region := config.Region.ValueString()
	awsopts = append(awsopts, awsconfig.WithRegion(region))
	ctx = tflog.SetField(ctx, "region", region)

	retryMaxAttempts := awsretry.DefaultMaxAttempts
	if !config.RetryMaxAttempts.IsNull() {
		retryMaxAttempts = int(config.RetryMaxAttempts.ValueInt64())
	}

	retryMaxBackoff := awsretry.DefaultMaxBackoff
	if !config.RetryMaxBackoff.IsNull() {
		retryMaxBackoff = time.Duration(config.RetryMaxBackoff.ValueInt64()) * time.Second
	}

	awsopts = append(awsopts, awsconfig.WithRetryer(func() aws2.Retryer {
		return awsretry.NewStandard(func(opts *awsretry.StandardOptions) {
			opts.MaxAttempts = retryMaxAttempts
			opts.MaxBackoff = retryMaxBackoff
		})
	}))
	tflog.Debug(ctx, "Using config", map[string]any{
		"region":       region,
		"max_attempts": retryMaxAttempts,
		"max_backoff":  retryMaxBackoff,
	})
	tflog.Debug(ctx, "Creating AWS AppStream client")

	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsopts...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create AWS config",
			"An unexpected error occurred when creating the AWS config. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"AWS Config Error: "+err.Error(),
		)
		return
	}

	clients := intaws.NewClients(awscfg)

	resp.DataSourceData = clients
	resp.ResourceData = clients

	tflog.Info(ctx, "Configured AWS AppStream client", map[string]any{"success": true})
}

func (A *awsAppStreamProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (A *awsAppStreamProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		entitlement.NewAssociateApplicationEntitlementResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &awsAppStreamProvider{version: version}
	}
}
