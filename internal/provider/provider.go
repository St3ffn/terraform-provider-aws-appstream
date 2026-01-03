// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsretry "github.com/aws/aws-sdk-go-v2/aws/retry"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/metadata"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/app_block"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/application"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/associate_application_entitlement"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/associate_application_fleet"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/associate_fleet_stack"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/associate_user_stack"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/directory_config"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/entitlement"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/fleet"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/stack"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/resources/user"
)

var (
	_ provider.Provider                   = &awsAppStreamProvider{}
	_ provider.ProviderWithValidateConfig = &awsAppStreamProvider{}
)

type awsAppStreamProvider struct {
	version string
}

// awsAppStreamProviderModel describes the provider data model.
type awsAppStreamProviderModel struct {
	AccessKey                 types.String      `tfsdk:"access_key"`
	SecretAccessKey           types.String      `tfsdk:"secret_access_key"`
	SessionToken              types.String      `tfsdk:"session_token"`
	Profile                   types.String      `tfsdk:"profile"`
	SkipCredentialsValidation types.Bool        `tfsdk:"skip_credentials_validation"`
	Region                    types.String      `tfsdk:"region"`
	RetryMode                 types.String      `tfsdk:"retry_mode"`
	RetryMaxAttempts          types.Int64       `tfsdk:"retry_max_attempts"`
	RetryMaxBackoff           types.Int64       `tfsdk:"retry_max_backoff"`
	DefaultTags               *defaultTagsModel `tfsdk:"default_tags"`
}

type defaultTagsModel struct {
	Tags types.Map `tfsdk:"tags"`
}

func (p *awsAppStreamProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.Version = p.version
	resp.TypeName = "awsappstream"
}

func (p *awsAppStreamProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with AWS AppStream",
		MarkdownDescription: `
The **awsappstream** provider allows you to manage AWS AppStream resources
using Terraform.

Authentication and region selection follow the standard AWS SDK behavior.
`,
		Attributes: map[string]schema.Attribute{
			"access_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Access Key ID. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: "The AWS access key ID to use for authentication. " +
					"If not set, the AWS SDK default credential resolution chain is used " +
					"(environment variables, shared credentials file, EC2/ECS metadata, etc.).",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"secret_access_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Access Key. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: "The AWS secret access key to use for authentication. " +
					"If not set, the AWS SDK default credential resolution chain is used " +
					"(environment variables, shared credentials file, EC2/ECS metadata, etc.).",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"session_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Session Token for temporary credentials. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: "The AWS session token to use for temporary credentials, such as those obtained via AWS STS. " +
					"This value is optional and typically only required when using temporary security credentials." +
					"If not set, the AWS SDK default credential resolution chain is used.",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "AWS Profile name. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: "The name of the AWS CLI profile to use. " +
					"If not set, the AWS SDK default credential resolution chain is used (environment variables, shared credentials file, EC2/ECS metadata, etc.).",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"skip_credentials_validation": schema.BoolAttribute{
				Optional:    true,
				Description: "Skip validation of AWS credentials via STS.",
				MarkdownDescription: "Skips validating AWS credentials using the STS `GetCallerIdentity` call. " +
					"Useful for testing or for AWS-compatible endpoints that do not support STS.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "AWS Region. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: "The AWS region in which AppStream resources are managed. " +
					"If not set, the AWS SDK default region resolution chain is used " +
					"(environment variables such as `AWS_REGION` or `AWS_DEFAULT_REGION`, " +
					"shared configuration files, or EC2/ECS metadata).",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"retry_mode": schema.StringAttribute{
				Optional:    true,
				Description: "Retry strategy used for AWS AppStream API calls. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: "Controls the retry strategy used by the AWS SDK when calling AWS AppStream APIs. " +
					"Supported values are:\n\n" +
					"\t- **`standard`** – Uses exponential backoff and retries failed requests *after* throttling or transient errors occur.\n" +
					"\t- **`adaptive`** – Dynamically adjusts the request rate based on AWS throttling responses, " +
					"slowing down *before* sending requests when throttling is detected. " +
					"This mode is recommended for workloads that create or update many AppStream resources concurrently.\n\n" +
					"\tIf not set, the AWS SDK default retry configuration is used (for example via environment variables such as `AWS_RETRY_MODE`).",
				Validators: []validator.String{
					stringvalidator.OneOf("standard", "adaptive"),
				},
			},
			"retry_max_attempts": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of attempts to retry AWS AppStream operations. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: fmt.Sprintf(
					"The maximum number of retry attempts for retryable AWS AppStream API requests. "+
						"Retries are only performed for retryable errors as determined by the AWS SDK "+
						"(for example throttling errors, transient network failures, and 5xx service errors). "+
						"Non-retryable errors such as validation or authorization failures are not retried. "+
						"If not set, the AWS SDK default retry configuration is used (for example via environment variables such as `AWS_MAX_ATTEMPTS`). "+
						"**SDK Default:** %d",
					awsretry.DefaultMaxAttempts,
				),
				Validators: []validator.Int64{int64validator.AtLeast(1)},
			},
			"retry_max_backoff": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum retry backoff in seconds for AWS AppStream operations. If unset, the AWS SDK default chain is used.",
				MarkdownDescription: fmt.Sprintf(
					"The maximum backoff time, in seconds, between retry attempts for retryable AWS AppStream API requests. "+
						"This limits the exponential backoff applied by the AWS SDK for retryable errors only. "+
						"If not set, the AWS SDK default retry configuration is used. "+
						"**SDK Default:** %d seconds",
					int64(awsretry.DefaultMaxBackoff.Seconds()),
				),
				Validators: []validator.Int64{int64validator.AtLeast(1)},
			},
			"default_tags": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Default tags to apply to all taggable resources managed by this provider.",
				MarkdownDescription: "Default tags to apply to all **taggable** resources managed by this provider. " +
					"Tags defined on individual resources take precedence over these defaults when keys overlap.",
				Attributes: map[string]schema.Attribute{
					"tags": schema.MapAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "A map of tags to apply by default.",
						MarkdownDescription: "A map of tags to apply by default. " +
							"Resource-level tags override these defaults when the same key is set.",
					},
				},
			},
		},
	}
}

func (p *awsAppStreamProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	tflog.Debug(ctx, "Validating AWS AppStream Provider")

	var config awsAppStreamProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.AccessKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Unknown AWS Access Key",
			"The AWS AppStream provider cannot be configured because \"access_key\" is unknown. "+
				"Provider configuration values must be static. "+
				"Set \"access_key\" to a fixed string or remove it to use the AWS SDK default.",
		)
		return
	}

	if config.SecretAccessKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_access_key"),
			"Unknown AWS Secret Access Key",
			"The AWS AppStream provider cannot be configured because \"secret_access_key\" is unknown. "+
				"Provider configuration values must be static. "+
				"Set \"secret_access_key\" to a fixed string or remove it to use the AWS SDK default.",
		)
		return
	}

	if config.SessionToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("session_token"),
			"Unknown AWS Session Token",
			"The AWS AppStream provider cannot be configured because \"session_token\" is unknown. "+
				"Provider configuration values must be static. "+
				"Set \"session_token\" to a fixed string or remove it to use the AWS SDK default.",
		)
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
				"Provider configuration values must be static. "+
				"Set \"region\" to a fixed string or remove it to use the AWS SDK default.",
		)
		return
	}

	if config.RetryMode.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("retry_mode"),
			"Unknown AWS Retry Mode",
			"The AWS AppStream provider cannot be configured because \"retry_mode\" is unknown. "+
				"Provider configuration values must be static. "+
				"Set \"retry_mode\" to a fixed value or remove it to use the default.",
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

	if config.DefaultTags != nil && config.DefaultTags.Tags.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("default_tags").AtName("tags"),
			"Unknown Default Tags",
			"The AWS AppStream provider cannot be configured because \"default_tags.tags\" is unknown. "+
				"Provider configuration values must be static. "+
				"Set \"default_tags.tags\" to a fixed map or remove it to use no default tags.",
		)
		return
	}

	hasAccessKey := !config.AccessKey.IsNull()
	hasSecretKey := !config.SecretAccessKey.IsNull()
	hasSession := !config.SessionToken.IsNull()
	hasProfile := !config.Profile.IsNull()

	if hasAccessKey != hasSecretKey {
		resp.Diagnostics.AddError(
			"Incomplete AWS Static Credentials",
			"The AWS AppStream provider cannot be configured because static credentials are incomplete. "+
				"Both \"access_key\" and \"secret_access_key\" must be set together, or both must be omitted.",
		)
		return
	}

	// session_token only makes sense with static credentials
	if hasSession && (!hasAccessKey && !hasSecretKey) {
		resp.Diagnostics.AddAttributeError(
			path.Root("session_token"),
			"Invalid AWS Session Token Configuration",
			"The AWS AppStream provider cannot be configured because \"session_token\" was set without static credentials. "+
				"Set \"access_key\" and \"secret_access_key\" together with \"session_token\", or remove \"session_token\" to use the AWS SDK default.",
		)
		return
	}

	if hasProfile && hasAccessKey {
		resp.Diagnostics.AddError(
			"Conflicting AWS Authentication Configuration",
			"The AWS AppStream provider cannot be configured because both a named AWS profile and static credentials were provided. "+
				"Set either \"profile\" or \"access_key\"/\"secret_access_key\", but not both.",
		)
		return
	}
}

func (p *awsAppStreamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring AWS AppStream Provider")

	var config awsAppStreamProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasAccessKey := !config.AccessKey.IsNull()
	hasSecretKey := !config.SecretAccessKey.IsNull()
	useStaticCreds := hasAccessKey && hasSecretKey
	hasSession := !config.SessionToken.IsNull()
	hasProfile := !config.Profile.IsNull()
	hasRegion := !config.Region.IsNull()
	hasRetryMode := !config.RetryMode.IsNull()
	hasRetryMaxAttempts := !config.RetryMaxAttempts.IsNull()
	hasRetryMaxBackoff := !config.RetryMaxBackoff.IsNull()
	retryConfigured := hasRetryMode || hasRetryMaxAttempts || hasRetryMaxBackoff

	var awsopts []func(*awsconfig.LoadOptions) error

	if useStaticCreds {
		creds := awscredentials.NewStaticCredentialsProvider(
			config.AccessKey.ValueString(),
			config.SecretAccessKey.ValueString(),
			func() string {
				if hasSession {
					return config.SessionToken.ValueString()
				}
				return ""
			}(),
		)
		awsopts = append(awsopts, awsconfig.WithCredentialsProvider(aws.NewCredentialsCache(creds)))
	} else if hasProfile {
		profile := config.Profile.ValueString()
		awsopts = append(awsopts, awsconfig.WithSharedConfigProfile(profile))
		ctx = tflog.SetField(ctx, "profile", profile)
	}

	switch {
	case useStaticCreds:
		tflog.Debug(ctx, "Using static AWS credentials from provider configuration")
	case hasProfile:
		tflog.Debug(ctx, "Using AWS profile from provider configuration")
	default:
		tflog.Debug(ctx, "Using AWS SDK default credential resolution chain")
	}

	if hasRegion {
		region := config.Region.ValueString()
		awsopts = append(awsopts, awsconfig.WithRegion(region))
		ctx = tflog.SetField(ctx, "region", region)
	}

	retryMaxAttempts := awsretry.DefaultMaxAttempts
	if hasRetryMaxAttempts {
		retryMaxAttempts = int(config.RetryMaxAttempts.ValueInt64())
	}

	retryMaxBackoff := awsretry.DefaultMaxBackoff
	if hasRetryMaxBackoff {
		retryMaxBackoff = time.Duration(config.RetryMaxBackoff.ValueInt64()) * time.Second
	}

	if retryConfigured {
		retryMode := "standard"
		if hasRetryMode {
			retryMode = config.RetryMode.ValueString()
		}

		awsopts = append(awsopts, awsconfig.WithRetryer(func() aws.Retryer {
			switch retryMode {
			case "adaptive":
				return awsretry.NewAdaptiveMode(func(o *awsretry.AdaptiveModeOptions) {
					o.StandardOptions = []func(*awsretry.StandardOptions){
						func(so *awsretry.StandardOptions) {
							so.MaxAttempts = retryMaxAttempts
							so.MaxBackoff = retryMaxBackoff
						},
					}
				})
			default:
				return awsretry.NewStandard(func(so *awsretry.StandardOptions) {
					so.MaxAttempts = retryMaxAttempts
					so.MaxBackoff = retryMaxBackoff
				})
			}
		}))

		tflog.Debug(ctx, "Using config", map[string]any{
			"retry_mode":         retryMode,
			"retry_max_attempts": retryMaxAttempts,
			"retry_max_backoff":  retryMaxBackoff,
		})
	}

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

	skipValidation := false
	if !config.SkipCredentialsValidation.IsNull() && !config.SkipCredentialsValidation.IsUnknown() {
		skipValidation = config.SkipCredentialsValidation.ValueBool()
	}

	if !skipValidation {
		stsClient := sts.NewFromConfig(awscfg)
		_, err = stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid AWS Credentials",
				fmt.Sprintf("Failed to validate AWS credentials: %v", err),
			)
			return
		}
	}

	defaultTags := map[string]string{}

	if config.DefaultTags != nil && !config.DefaultTags.Tags.IsNull() {
		diags := config.DefaultTags.Tags.ElementsAs(ctx, &defaultTags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	meta := metadata.NewMetadata(awscfg, defaultTags)

	resp.DataSourceData = meta
	resp.ResourceData = meta

	tflog.Info(ctx, "Configured AWS AppStream client", map[string]any{"success": true})
}

func (p *awsAppStreamProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		fleet.NewDataSource,
		stack.NewDataSource,
		entitlement.NewDataSource,
		app_block.NewDataSource,
		application.NewDataSource,
		directory_config.NewDataSource,
		user.NewDataSource,
	}
}

func (p *awsAppStreamProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		fleet.NewResource,
		stack.NewResource,
		entitlement.NewResource,
		app_block.NewResource,
		application.NewResource,
		directory_config.NewResource,
		user.NewResource,
		associate_fleet_stack.NewResource,
		associate_application_entitlement.NewResource,
		associate_application_fleet.NewResource,
		associate_user_stack.NewResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &awsAppStreamProvider{version: version}
	}
}
