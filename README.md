# Terraform AWS AppStream Provider

[![Release](https://img.shields.io/github/v/release/St3ffn/terraform-provider-aws-appstream)](https://github.com/St3ffn/terraform-provider-aws-appstream/releases)
[![CI](https://github.com/St3ffn/terraform-provider-aws-appstream/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/St3ffn/terraform-provider-aws-appstream/actions/workflows/test.yml?branch=main)
[![License](https://img.shields.io/github/license/st3ffn/terraform-provider-aws-appstream)](/LICENSE)
[![GO](https://img.shields.io/github/go-mod/go-version/St3ffn/terraform-provider-aws-appstream)](https://golang.org/)

---

The **awsappstream** provider allows Terraform to manage AWS AppStream resources.

Authentication and region configuration follow the standard AWS SDK behavior,
with optional customization of retry behavior.

## Motivation

Support for AWS AppStream in the official
[`hashicorp/terraform-provider-aws`](https://github.com/hashicorp/terraform-provider-aws) provider is currently 
incomplete and, in some cases, does not reflect the latest AWS AppStream features or expected resource behavior.

This provider was created to offer more complete and reliable AppStream support.
For context, [see the open AppStream-related issues](https://github.com/hashicorp/terraform-provider-aws/issues?q=is%3Aissue%20state%3Aopen%20label%3Aservice%2Fappstream) in the AWS provider.

## AppStream Coverage

| Name                                           | Resource | Data Source | Planned |
|------------------------------------------------|----------|-------------|---------|
| awsappstream_fleet                             | ✅        | ✅           |         |
| awsappstream_stack                             | ✅        | ✅           |         |
| awsappstream_associate_fleet_stack             | ✅        | ❌           |         |
| awsappstream_entitlement                       | ✅        | ✅           |         |
| awsappstream_associate_application_entitlement | ✅        | ❌           |         |
| awsappstream_app_block                         | ✅        | ✅           |         |
| awsappstream_application                       | ✅        | ✅           |         |
| awsappstream_directory_config                  | ✅        | ✅           |         |
| awsappstream_user                              | ✅        | ✅           |         |
| awsappstream_associate_user_stack              | ✅        | ❌           |         |
| awsappstream_associate_application_fleet       | ✅        | ❌           |         |
| awsappstream_image                             | ❌        | ✅           |         |
| awsappstream_image_builder                     | 🚧       | 🚧          | ✅       |

## Behavior and Design Principles

This provider follows a **read-after-write** model to ensure Terraform state
accurately reflects the authoritative state in AWS.

In practice, this means:

- **Create and Update operations are always followed by a Read**
    - After a successful `Create` or `Update`, the provider performs a fresh read
      from AWS and uses that response as the source of truth for state.
    - This avoids relying on partial or inconsistent API responses.

- **Read is authoritative**
    - If a resource cannot be found during `Read`, it is removed from state.
    - This applies to external deletions and drift detection.

- **Idempotent behavior**
    - Create operations tolerate existing resources where possible and converge
      state instead of failing when safe to do so.
    - Association-style resources model relationships only and verify existence
      rather than storing mutable state.

- **Tag management is declarative**
    - Tags are reconciled using a diff-based approach.
    - Default tags and resource-level tags are merged and applied consistently
      during Create and Update.
    - Changes to default tags are automatically propagated on the next apply.

- **Context-aware cancellation**
    - All operations respect context cancellation and deadlines to avoid
      corrupting state during interrupted applies.

Overall, the provider prioritizes **correctness, consistency, and drift
resilience** over minimizing API calls.

## Provider Design Philosophy: Attribute Ownership

This provider intentionally follows an **attribute ownership model** that differs from the official Terraform AWS provider.

Only attributes that are **explicitly set by the user in Terraform configuration** are considered managed by Terraform. Attributes that are not configured by the user are treated as **AWS-managed defaults** and are not enforced or reconciled by Terraform.

In practical terms:

- If an attribute is **not set** in the Terraform configuration, the provider:
    - Does not attempt to normalize it
    - Does not enforce AWS default values
    - Does not generate diffs when AWS populates or changes default values
- If an attribute **is set** by the user, Terraform fully owns it:
    - The value is sent to AWS
    - The value is read back from AWS
    - Drift is detected and corrected if the value changes

This behavior avoids perpetual diffs caused by AWS-side defaults, implicit behavior, or service-specific normalization.

### Examples

Fleet resources:
- If `image_name` is set and `image_arn` is not, only `image_name` is tracked.
- If `image_arn` is set, only `image_arn` is tracked.
- If AWS returns both values, the provider preserves ownership of only the attribute configured by the user.

Stack resources:
- Optional nested blocks (for example, user settings, application settings, streaming experience settings) are only tracked if explicitly defined.
- AWS-generated defaults or inferred values are ignored unless the user opts in by defining them.

### Why this design exists

AWS AppStream frequently:
- Applies implicit defaults
- Normalizes values
- Returns additional fields that were never explicitly configured

Tracking all returned fields would lead to:
- Constant, non-actionable diffs
- Forced configuration of values users did not intend to manage
- Reduced clarity around which settings Terraform truly controls

By limiting state ownership to user-defined attributes, this provider:
- Produces stable plans
- Avoids perpetual drift
- Makes ownership boundaries explicit and predictable

### Implications for users

- You are not required to configure every available attribute.
- AWS defaults are respected unless you choose to override them.
- If you want Terraform to manage a value, you must explicitly define it.
- Importing existing resources will populate only user-managed attributes.

This behavior is intentional and applies consistently across resources such as fleets and stacks, and may be extended to additional resources in the future.

## Retry and Eventual Consistency Handling

AWS AppStream APIs exhibit **eventual consistency** and transient errors,
especially when creating or associating dependent resources (for example,
fleets, stacks, and entitlements).

This provider uses a **layered retry approach**:

- **AWS SDK retries**
    - Configurable via provider settings (`retry_mode`, `retry_max_attempts`, `retry_max_backoff`)
    - Handles throttling, networking issues, and standard AWS retryable errors

- **Provider-level retries**
    - Applied selectively to operations known to fail temporarily due to
      AppStream lifecycle constraints (for example `OperationNotPermittedException`
      or `ResourceNotFoundException` during creation or association)
    - Uses bounded exponential backoff and respects Terraform cancellation

This ensures Terraform operations converge reliably without requiring
manual sleeps or explicit dependencies in configuration.
