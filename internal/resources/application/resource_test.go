// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccApplicationBasicConfig(appBlockName, applicationName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %q

  packaging_type = "CUSTOM"

  source_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = "appstream-acc-test-bucket"
      s3_key    = "app_block_setup.sh"
    }

    executable_path    = "/bin/bash"
    timeout_in_seconds = 60
  }
}

resource "awsappstream_application" "test" {
  name = "%s"

  icon_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  launch_path = "/app/app"

  platforms = ["AMAZON_LINUX2"]
  instance_families = ["GENERAL_PURPOSE"]

  app_block_arn = awsappstream_app_block.test.arn
}
`, appBlockName, applicationName)
}

func TestAccApplication_basic(t *testing.T) {
	appBlockName := acctest.RandomWithPrefix("tf-acc-app-block")
	applicationName := acctest.RandomWithPrefix("tf-acc-application")

	resourceName := "awsappstream_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationBasicConfig(appBlockName, applicationName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", applicationName),
					resource.TestCheckResourceAttr(resourceName, "launch_path", "/app/app"),
					resource.TestCheckResourceAttr(resourceName, "platforms.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_families.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "app_block_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckNoResourceAttr(resourceName, "tags"),
					resource.TestCheckNoResourceAttr(resourceName, "tags_all"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccApplicationBasicConfig(appBlockName, applicationName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccApplicationComplexConfig(appBlockName, applicationName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %q

  packaging_type = "CUSTOM"

  source_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = "appstream-acc-test-bucket"
      s3_key    = "app_block_setup.sh"
    }

    executable_path    = "/bin/bash"
    timeout_in_seconds = 60
  }
}

resource "awsappstream_application" "test" {
  name = "%s"

  display_name = "Test Application"
  description  = "Complex acceptance test"

  icon_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  launch_path        = "/app/app"
  working_directory  = "/app"
  launch_parameters  = "--debug"

  platforms = ["AMAZON_LINUX2"]
  instance_families = ["GENERAL_PURPOSE"]

  app_block_arn = awsappstream_app_block.test.arn

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`, appBlockName, applicationName)
}

func testAccApplicationComplexConfigUpdated(appBlockName, applicationName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %q

  packaging_type = "CUSTOM"

  source_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = "appstream-acc-test-bucket"
      s3_key    = "app_block_setup.sh"
    }

    executable_path    = "/bin/bash"
    timeout_in_seconds = 60
  }
}

resource "awsappstream_application" "test" {
  name = "%s"

  display_name = "Updated Application"
  description  = "Updated complex acceptance test"

  icon_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  launch_path        = "/app/app"
  working_directory  = "/app/bin"
  launch_parameters  = "--verbose"

  platforms = ["AMAZON_LINUX2"]
  instance_families = ["GENERAL_PURPOSE"]

  app_block_arn = awsappstream_app_block.test.arn

  tags = {
    Environment = "prod"
    Team        = "platform"
  }
}
`, appBlockName, applicationName)
}

func TestAccApplication_complex(t *testing.T) {
	appBlockName := acctest.RandomWithPrefix("tf-acc-app-block")
	applicationName := acctest.RandomWithPrefix("tf-acc-application")

	resourceName := "awsappstream_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationComplexConfig(appBlockName, applicationName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Application"),
					resource.TestCheckResourceAttr(resourceName, "description", "Complex acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "working_directory", "/app"),
					resource.TestCheckResourceAttr(resourceName, "launch_parameters", "--debug"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Owner", "terraform"),
				),
			},
			{
				Config: testAccApplicationComplexConfigUpdated(appBlockName, applicationName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Application"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated complex acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "working_directory", "/app/bin"),
					resource.TestCheckResourceAttr(resourceName, "launch_parameters", "--verbose"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "prod"),
					resource.TestCheckResourceAttr(resourceName, "tags.Team", "platform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Environment", "prod"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Team", "platform"),
				),
			},
		},
	})
}

func testAccApplicationTagsConfig(appBlockName, applicationName string) string {
	return testhelpers.TestAccProviderTagsConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %q

  packaging_type = "CUSTOM"

  source_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = "appstream-acc-test-bucket"
      s3_key    = "app_block_setup.sh"
    }

    executable_path    = "/bin/bash"
    timeout_in_seconds = 60
  }
}

resource "awsappstream_application" "test" {
  name = "%s"

  icon_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  launch_path = "/app/app"

  platforms = ["AMAZON_LINUX2"]
  instance_families = ["GENERAL_PURPOSE"]

  app_block_arn = awsappstream_app_block.test.arn

  tags = {
    Environment = "test"
  }
}
`, appBlockName, applicationName)
}

func TestAccApplication_tags(t *testing.T) {
	appBlockName := acctest.RandomWithPrefix("tf-acc-app-block")
	applicationName := acctest.RandomWithPrefix("tf-acc-application")

	resourceName := "awsappstream_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationTagsConfig(appBlockName, applicationName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.MANAGED_BY", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.BUILD_WITH", "love"),
				),
			},
		},
	})
}
