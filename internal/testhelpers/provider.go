// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/provider"
)

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"awsappstream": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// TestAccProviderBasicConfig returns a minimal provider configuration
// without any custom settings or default tags.
func TestAccProviderBasicConfig() string {
	return `
provider "awsappstream" {}
`
}

// TestAccProviderTagsConfig returns a provider configuration with default_tags
// enabled for acceptance tests.
//
// The following default tags are defined:
//   - MANAGED_BY = "terraform"
//   - BUILD_WITH = "love"
func TestAccProviderTagsConfig() string {
	return `
provider "awsappstream" {
  default_tags {
    tags = {
      MANAGED_BY  = "terraform"
      BUILD_WITH  = "love"
    }
  }
}
`
}
