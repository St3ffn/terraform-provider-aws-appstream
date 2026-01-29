// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccAppBlockBuilderWithDataSource(name string, subnetIDs []string) string {
	return testhelpers.TestAccProviderTagsConfig() + fmt.Sprintf(`
resource "awsappstream_app_block_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  platform      = "WINDOWS_SERVER_2019"

  vpc_config = {
    subnet_ids = [%s]
  }

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}

data "awsappstream_app_block_builder" "test" {
  name = awsappstream_app_block_builder.test.name

  depends_on = [awsappstream_app_block_builder.test]
}
`, name, testhelpers.HCLStringList(subnetIDs))
}

func TestAccAppBlockBuilderDataSource_basic(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)
	vpc := testCtx.DefaultVPCInfo(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-builder-ds-basic")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBuilderWithDataSource(name, vpc.SubnetIDs[:2]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_app_block_builder.test", "name", name),
					resource.TestCheckResourceAttr("data.awsappstream_app_block_builder.test", "instance_type", "stream.standard.small"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block_builder.test", "platform", "WINDOWS_SERVER_2019"),
					resource.TestCheckResourceAttrSet("data.awsappstream_app_block_builder.test", "enable_default_internet_access"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block_builder.test", "vpc_config.subnet_ids.#", "2"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block_builder.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block_builder.test", "tags.Owner", "terraform"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block_builder.test", "tags.MANAGED_BY", "terraform"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block_builder.test", "tags.BUILD_WITH", "love"),
					resource.TestCheckResourceAttrSet("data.awsappstream_app_block_builder.test", "arn"),
					resource.TestCheckResourceAttrSet("data.awsappstream_app_block_builder.test", "created_time"),
					resource.TestCheckResourceAttrSet("data.awsappstream_app_block_builder.test", "state"),
				),
			},
		},
	})
}
