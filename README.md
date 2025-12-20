# Terraform AWS AppStream Provider

The **awsappstream** provider allows Terraform to manage AWS AppStream resources.

Authentication and region configuration follow the standard AWS SDK behavior,
with optional customization of retry behavior.

## Motivation

Support for AWS AppStream in the official
[`hashicorp/terraform-provider-aws`](https://github.com/hashicorp/terraform-provider-aws) provider is currently 
incomplete and, in some cases, does not reflect the latest AWS AppStream features or expected resource behavior.

This provider was created to offer more complete and reliable AppStream support.
For context, [see the open AppStream-related issues](https://github.com/hashicorp/terraform-provider-aws/issues?q=is%3Aissue%20state%3Aopen%20label%3Aservice%2Fappstream) in the AWS provider.

## Supported Resources

- `awsappstream_entitlement`
- `awsappstream_associate_application_entitlement`

## Supported Data Sources

- `awsappstream_stack`
- `awsappstream_entitlement`
