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

func TestAccProviderBasicConfig() string {
	return `
provider "awsappstream" {}
`
}
