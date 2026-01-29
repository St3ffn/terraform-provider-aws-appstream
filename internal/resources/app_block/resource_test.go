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

func testAccAppBlockBasicConfig(name, bucketName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
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
}
`, name, bucketName)
}

func TestAccAppBlock_basic(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block")
	resourceName := "awsappstream_app_block.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBasicConfig(name, testCtx.BucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "packaging_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "source_s3_location.s3_bucket", testCtx.BucketName),
					resource.TestCheckResourceAttr(resourceName, "source_s3_location.s3_key", "appblock.vhd"),
					resource.TestCheckResourceAttr(resourceName, "setup_script_details.executable_path", "/bin/bash"),
					resource.TestCheckResourceAttr(resourceName, "setup_script_details.timeout_in_seconds", "60"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckNoResourceAttr(resourceName, "tags"),
					resource.TestCheckNoResourceAttr(resourceName, "tags_all"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setup_script_details",
				},
			},
		},
	})
}

func testAccAppBlockComplexConfig(name, bucketName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %[1]q

  packaging_type = "CUSTOM"

  display_name = "Test App Block"
  description  = "Complex acceptance test"

  source_s3_location = {
    s3_bucket = %[2]q
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = %[2]q
      s3_key    = "app_block_setup.sh"
    }

    executable_path       = "/bin/bash"
    executable_parameters = "-e"
    timeout_in_seconds    = 10
  }

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`, name, bucketName)
}

func testAccAppBlockComplexConfigUpdated(name, bucketName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %[1]q

  packaging_type = "CUSTOM"

  display_name = "Updated App Block"
  description  = "Updated complex acceptance test"

  source_s3_location = {
    s3_bucket = %[2]q
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = %[2]q
      s3_key    = "app_block_setup.sh"
    }

    executable_path       = "/bin/bash"
    executable_parameters = "-x"
    timeout_in_seconds    = 60
  }

  tags = {
    Environment = "prod"
    Team        = "platform"
  }
}
`, name, bucketName)
}

func TestAccAppBlock_complex(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-complex")
	resourceName := "awsappstream_app_block.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockComplexConfig(name, testCtx.BucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test App Block"),
					resource.TestCheckResourceAttr(resourceName, "description", "Complex acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "packaging_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "source_s3_location.s3_bucket", testCtx.BucketName),
					resource.TestCheckResourceAttr(resourceName, "source_s3_location.s3_key", "appblock.vhd"),
					resource.TestCheckResourceAttr(resourceName, "setup_script_details.executable_path", "/bin/bash"),
					resource.TestCheckResourceAttr(resourceName, "setup_script_details.executable_parameters", "-e"),
					resource.TestCheckResourceAttr(resourceName, "setup_script_details.timeout_in_seconds", "10"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			{
				Config: testAccAppBlockComplexConfigUpdated(name, testCtx.BucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated App Block"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated complex acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "setup_script_details.executable_parameters", "-x"),
					resource.TestCheckResourceAttr(resourceName, "setup_script_details.timeout_in_seconds", "60"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "prod"),
					resource.TestCheckResourceAttr(resourceName, "tags.Team", "platform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Environment", "prod"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Team", "platform"),
				),
			},
		},
	})
}

func TestAccAppBlock_noopPlan(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-noop")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBasicConfig(name, testCtx.BucketName),
			},
			{
				Config:             testAccAppBlockBasicConfig(name, testCtx.BucketName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccAppBlockAppStream2Config(name, bucketName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %[1]q

  packaging_type = "APPSTREAM2"

  display_name = "AppStream2 App Block"
  description  = "APPSTREAM2 acceptance test"

  source_s3_location = {
    s3_bucket = %[2]q
  }

  post_setup_script_details = {
    script_s3_location = {
      s3_bucket = %[2]q
      s3_key    = "app_block_post_setup.sh"
    }

    executable_path    = "/bin/bash"
    timeout_in_seconds = 30
  }

  tags = {
    Environment = "test"
  }
}
`, name, bucketName)
}

func TestAccAppBlock_appstream2(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-appstream2")
	resourceName := "awsappstream_app_block.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockAppStream2Config(name, testCtx.BucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "packaging_type", "APPSTREAM2"),
					resource.TestCheckResourceAttr(resourceName, "post_setup_script_details.executable_path", "/bin/bash"),
					resource.TestCheckResourceAttr(resourceName, "post_setup_script_details.timeout_in_seconds", "30"),
					resource.TestCheckNoResourceAttr(resourceName, "setup_script_details"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Environment", "test"),
				),
			},
		},
	})
}

func testAccAppBlockTagsConfig(name, bucketName string) string {
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
  }
}
`, name, bucketName)
}

func TestAccAppBlock_tags(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-tags")
	resourceName := "awsappstream_app_block.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockTagsConfig(name, testCtx.BucketName),
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
