// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccStackWithDataSource(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

data "awsappstream_stack" "test" {
  name = awsappstream_stack.test.name
}
`, name)
}

func TestAccStackDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-ds-basic")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackWithDataSource(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_stack.test", "name", name),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack.test", "arn"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack.test", "created_time"),
				),
			},
		},
	})
}

func testAccStackWithDataSourceComplex(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q

  redirect_url = "https://example.com/logout"

  application_settings = {
    enabled = false
  }

  user_settings = [
    {
      action     = "CLIPBOARD_COPY_FROM_LOCAL_DEVICE"
      permission = "ENABLED"
    }
  ]
}

data "awsappstream_stack" "test" {
  name = awsappstream_stack.test.name
}
`, name)
}

func TestAccStackDataSource_complex(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-ds-complex")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackWithDataSourceComplex(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_stack.test", "name", name),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack.test", "arn"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack.test", "created_time"),
					resource.TestCheckResourceAttr("data.awsappstream_stack.test", "application_settings.enabled", "false"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack.test", "user_settings.#"),
				),
			},
		},
	})
}
