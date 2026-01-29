// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccAppBlockWithDataSource(name, bucketName string) string {
	return testhelpers.TestAccProviderTagsConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %[1]q

  source_s3_location = {
    s3_bucket = %[2]q
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = %[2]q
      s3_key    = "app_block_setup.sh"
    }

    executable_path    = "/bin/bash"
    timeout_in_seconds = 60
  }

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}

data "awsappstream_app_block" "test" {
  arn = awsappstream_app_block.test.arn
  
  depends_on = [awsappstream_app_block.test]
}
`, name, bucketName)
}

func TestAccAppBlockDataSource_basic(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-ds-basic")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockWithDataSource(name, testCtx.BucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_app_block.test", "name", name),
					resource.TestCheckResourceAttr("data.awsappstream_app_block.test", "packaging_type", "CUSTOM"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block.test", "source_s3_location.s3_bucket", testCtx.BucketName),
					resource.TestCheckResourceAttr("data.awsappstream_app_block.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block.test", "tags.Owner", "terraform"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block.test", "tags.MANAGED_BY", "terraform"),
					resource.TestCheckResourceAttr("data.awsappstream_app_block.test", "tags.BUILD_WITH", "love"),
					resource.TestCheckResourceAttrSet("data.awsappstream_app_block.test", "arn"),
					resource.TestCheckResourceAttrSet("data.awsappstream_app_block.test", "created_time"),
					resource.TestCheckResourceAttrSet("data.awsappstream_app_block.test", "state"),
				),
			},
		},
	})
}
