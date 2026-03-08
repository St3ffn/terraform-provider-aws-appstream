---
page_title: "Migration Guide: hashicorp/aws to st3ffn/aws-appstream"
subcategory: "Migration"
description: |-
  Guidance for migrating AppStream configurations from the hashicorp/aws provider to the st3ffn/awsappstream provider.
---

# Migration Guide: `hashicorp/aws` to `st3ffn/aws-appstream`

This guide highlights the main differences when moving AppStream resources from `hashicorp/aws` to `st3ffn/aws-appstream`.
It is intended as a practical checklist for existing Terraform configurations.

## Why Migrate

- The `aws-appstream` provider focuses only on AppStream resources.
- AppStream-specific lifecycle behavior is implemented explicitly (for example fleet stop/start update handling).
- AppStream schemas and provider behavior are aligned and documented in one place.

## Provider Configuration

Configure both providers separately during migration:

```terraform
terraform {
  required_version = ">= 1.2"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
    awsappstream = {
      source  = "st3ffn/aws-appstream"
      version = "~> 2.0"
    }
  }
}

provider "aws" {
  region = "eu-central-1"
}

provider "awsappstream" {
  region = "eu-central-1"
}
```

Use `aws` for non-AppStream dependencies (VPC, IAM, SGs) and `awsappstream` for AppStream resources.

## Resource Type Naming

The provider type prefix differs:

- `hashicorp/aws`: `aws_appstream_*`
- `st3ffn/aws-appstream`: `awsappstream_*`

Example:

- `aws_appstream_fleet` -> `awsappstream_fleet`
- `aws_appstream_stack` -> `awsappstream_stack`

Association resources may also use different naming conventions (for example verb-first names such as `awsappstream_associate_fleet_stack`).

## Resource Mapping

The table below maps AppStream **resources** from `hashicorp/aws` to this provider.
Baseline used here: `hashicorp/aws` `v6.34.0`.

| `hashicorp/aws` resource                | `st3ffn/aws-appstream` resource      |
|-----------------------------------------|--------------------------------------|
| `aws_appstream_directory_config`        | `awsappstream_directory_config`      |
| `aws_appstream_fleet`                   | `awsappstream_fleet`                 |
| `aws_appstream_fleet_stack_association` | `awsappstream_associate_fleet_stack` |
| `aws_appstream_image_builder`           | `awsappstream_image_builder`         |
| `aws_appstream_stack`                   | `awsappstream_stack`                 |
| `aws_appstream_user`                    | `awsappstream_user`                  |
| `aws_appstream_user_stack_association`  | `awsappstream_associate_user_stack`  |

Resources currently available only in `st3ffn/aws-appstream`:

- `awsappstream_app_block`
- `awsappstream_app_block_builder`
- `awsappstream_application`
- `awsappstream_associate_app_block_builder_app_block`
- `awsappstream_associate_application_entitlement`
- `awsappstream_associate_application_fleet`
- `awsappstream_associate_image_builder_software`
- `awsappstream_entitlement`
- `awsappstream_copied_image`
- `awsappstream_imported_image`
- `awsappstream_updated_image`
- `awsappstream_image_permission`
- `awsappstream_stack_theme`
- `awsappstream_usage_report_subscription`

## Syntax Differences

`st3ffn/aws-appstream` schemas model nested values as attributes.
Use attribute assignment (`=`), not nested block syntax, for nested objects/collections.

Example:

```terraform
resource "awsappstream_fleet" "example" {
  name          = "example-fleet"
  fleet_type    = "ON_DEMAND"
  image_name    = "example-image"
  instance_type = "stream.standard.small"

  compute_capacity = {
    desired_instances = 1
  }

  vpc_config = {
    subnet_ids         = ["subnet-123", "subnet-456"]
    security_group_ids = ["sg-123"]
  }
}
```

For collections of objects, use list/set values:

```terraform
resource "awsappstream_stack" "example" {
  name = "example-stack"

  storage_connectors = [
    {
      connector_type = "HOMEFOLDERS"
    },
    {
      connector_type = "ONE_DRIVE"
      domains        = ["example.com"]
    }
  ]
}
```

Do not mix `storage_connectors = ...` with `dynamic "storage_connectors"` block syntax for these schemas.

## `lifecycle.ignore_changes` Paths

For nested attributes configured as objects, use dot paths:

```terraform
lifecycle {
  ignore_changes = [
    compute_capacity.desired_instances,
  ]
}
```

Do not use list index paths (for example `compute_capacity[0].desired_instances`) for object-typed nested attributes.

## Runtime/Lifecycle Behavior

Fleet updates can require stop/start behavior. Use:

- `update_behavior` (`AUTO_STOP_START` or `FAIL_IF_RUNNING`)
- `desired_state` (`INHERIT`, `RUNNING`, `STOPPED`)

Plan-time warnings may appear when a fleet update is expected to require restart behavior.

## Supported Resources and Data Sources

At a high level, `st3ffn/awsappstream` includes:

- Core resources: `fleet`, `stack`, `image_builder`, `app_block_builder`, `application`, `entitlement`, `directory_config`, `user`, `usage_report_subscription`, `stack_theme`, `copied_image`, `imported_image`, `updated_image`, `image_permission`, `app_block`
- Association resources: `associate_fleet_stack`, `associate_application_fleet`, `associate_application_entitlement`, `associate_user_stack`, `associate_image_builder_software`, `associate_app_block_builder_app_block`
- Data sources for core entities plus operational views such as `sessions`, `software_associations`, `export_image_task`, and `export_image_tasks`

Check the resource/data source docs to confirm exact fields and constraints before cutover.

## Write-Only/Sensitive Fields

Some credentials are modeled as write-only where supported (for example directory config account passwords).
These values are accepted from configuration but are not stored in Terraform state.

## Migration Strategy

1. Import existing AppStream resources into `awsappstream`.
2. Compare plans and fix syntax differences (`block` vs `attribute` assignment).
3. Migrate core resources first (fleet/stack/image-builder/app-block-builder), then association resources.
4. Remove old `aws` AppStream resource definitions after state is stable.
