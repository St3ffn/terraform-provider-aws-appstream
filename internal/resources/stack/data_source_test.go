// Copyright St3ffn 2025, 2026
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

  depends_on = [awsappstream_stack.test]
}
`, name)
}

func TestAccStackDataSource_basic(t *testing.T) {
	testhelpers.LoadAccTestEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-stack-ds-basic")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackWithDataSource(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_stack.test", "name", name),
					resource.TestCheckNoResourceAttr("data.awsappstream_stack.test", "tags"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack.test", "arn"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack.test", "created_time"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack.test", "created_time"),
				),
			},
		},
	})
}

func testAccStackWithDataSourceComplex(name string) string {
	return testhelpers.TestAccProviderTagsConfig() + fmt.Sprintf(`
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

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}

data "awsappstream_stack" "test" {
  name = awsappstream_stack.test.name

  depends_on = [awsappstream_stack.test]
}
`, name)
}

func TestAccStackDataSource_complex(t *testing.T) {
	testhelpers.LoadAccTestEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-stack-ds-complex")

	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr("data.awsappstream_stack.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("data.awsappstream_stack.test", "tags.Owner", "terraform"),
					resource.TestCheckResourceAttr("data.awsappstream_stack.test", "tags.MANAGED_BY", "terraform"),
					resource.TestCheckResourceAttr("data.awsappstream_stack.test", "tags.BUILD_WITH", "love"),
				),
			},
		},
	})
}
