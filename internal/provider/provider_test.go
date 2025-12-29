// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func TestAccProvider_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		PreCheck: func() {
			testhelpers.TestAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testhelpers.TestAccProviderBasicConfig(),
			},
		},
	})
}
