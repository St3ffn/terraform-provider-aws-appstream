// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build generate

package tools

import (
	_ "github.com/hashicorp/copywrite"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)

// Generate schema-derived artifacts (`*_model_gen.go` and `*_diff_gen.go`).
//go:generate go run ./cmd/schema-to-artifacts --root ../../internal/resources --schema-pattern resource_schema.go --schema-pattern data_source_schema.go

// Generate license.
//go:generate go run github.com/hashicorp/copywrite license -d ../.. --config ../../.copywrite.hcl

// Generate/update copyright headers.
//go:generate go run github.com/hashicorp/copywrite headers -d ../.. --config ../../.copywrite.hcl

// Format Terraform examples used in docs.
//go:generate terraform fmt -recursive ../../examples/

// Generate provider documentation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name awsappstream --rendered-provider-name "AWS AppStream" --provider-dir ../..
