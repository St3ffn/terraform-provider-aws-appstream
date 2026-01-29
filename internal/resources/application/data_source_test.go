// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package application_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccApplicationWithDataSource(appBlockName, applicationName, bucketName string) string {
	return testhelpers.TestAccProviderTagsConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %[1]q

  packaging_type = "CUSTOM"

  source_s3_location = {
    s3_bucket = %[3]q
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = %[3]q
      s3_key    = "app_block_setup.sh"
    }

    executable_path    = "/bin/bash"
    timeout_in_seconds = 60
  }
}

resource "awsappstream_application" "test" {
  name = "%[2]s"

  display_name = "Test Application"
  description  = "Application data source test"

  icon_s3_location = {
    s3_bucket = %[3]q
    s3_key    = "application_icon.png"
  }

  launch_path       = "/app/app"
  working_directory = "/app"
  launch_parameters = "--debug"

  platforms = ["AMAZON_LINUX2"]
  instance_families = ["GENERAL_PURPOSE"]

  app_block_arn = awsappstream_app_block.test.arn

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}

data "awsappstream_application" "test" {
  arn = awsappstream_application.test.arn

  depends_on = [awsappstream_application.test]
}
`, appBlockName, applicationName, bucketName)
}

func TestAccApplicationDataSource_basic(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)

	appBlockName := acctest.RandomWithPrefix("tf-acc-app-block-ds")
	applicationName := acctest.RandomWithPrefix("tf-acc-application-ds")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationWithDataSource(appBlockName, applicationName, testCtx.BucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "name", applicationName),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "display_name", "Test Application"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "description", "Application data source test"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "launch_path", "/app/app"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "working_directory", "/app"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "launch_parameters", "--debug"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "platforms.#", "1"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "instance_families.#", "1"),
					resource.TestCheckResourceAttrSet("data.awsappstream_application.test", "app_block_arn"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "tags.Owner", "terraform"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "tags.MANAGED_BY", "terraform"),
					resource.TestCheckResourceAttr("data.awsappstream_application.test", "tags.BUILD_WITH", "love"),
					resource.TestCheckResourceAttrSet("data.awsappstream_application.test", "arn"),
					resource.TestCheckResourceAttrSet("data.awsappstream_application.test", "id"),
					resource.TestCheckResourceAttrSet("data.awsappstream_application.test", "created_time"),
				),
			},
		},
	})
}
