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
| awsappstream_fleet_stack_association           | ✅        | ✅           |         |
| awsappstream_entitlement                       | ✅        | ✅           |         |
| awsappstream_associate_application_entitlement | ✅        | ✅           |         |

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

