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

func testAccAppBlockBuilderBasicConfig(name string, subnetIDs []string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  platform      = "WINDOWS_SERVER_2019"

  vpc_config = {
    subnet_ids = [%s]
  }
}
`, name, testhelpers.HCLStringList(subnetIDs))
}

func TestAccAppBlockBuilder_basic(t *testing.T) {
	testEnv := testhelpers.LoadAccTestEnv(t)
	vpc := testEnv.DefaultVPCInfo(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-builder")
	resourceName := "awsappstream_app_block_builder.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBuilderBasicConfig(name, vpc.SubnetIDs[:2]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "stream.standard.small"),
					resource.TestCheckResourceAttr(resourceName, "platform", "WINDOWS_SERVER_2019"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_default_internet_access"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.subnet_ids.#", "2"),
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
			},
		},
	})
}

func testAccAppBlockBuilderComplexConfig(name string, subnetIDs []string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  platform      = "WINDOWS_SERVER_2019"

  display_name = "Test App Block Builder"
  description  = "Complex acceptance test"

  vpc_config = {
    subnet_ids = [%s]
  }

  enable_default_internet_access = true

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`, name, testhelpers.HCLStringList(subnetIDs))
}

func testAccAppBlockBuilderComplexConfigUpdated(name string, subnetIDs []string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  platform      = "WINDOWS_SERVER_2019"

  display_name = "Updated App Block Builder"
  description  = "Updated complex acceptance test"

  vpc_config = {
    subnet_ids = [%s]
  }

  enable_default_internet_access = false

  tags = {
    Environment = "prod"
    Team        = "platform"
  }
}
`, name, testhelpers.HCLStringList(subnetIDs))
}

func TestAccAppBlockBuilder_complex(t *testing.T) {
	testEnv := testhelpers.LoadAccTestEnv(t)
	vpc := testEnv.DefaultVPCInfo(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-builder")
	resourceName := "awsappstream_app_block_builder.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBuilderComplexConfig(name, vpc.SubnetIDs[:2]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test App Block Builder"),
					resource.TestCheckResourceAttr(resourceName, "description", "Complex acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "stream.standard.small"),
					resource.TestCheckResourceAttr(resourceName, "platform", "WINDOWS_SERVER_2019"),
					resource.TestCheckResourceAttr(resourceName, "enable_default_internet_access", "true"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.subnet_ids.#", "2"),
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
				Config: testAccAppBlockBuilderComplexConfigUpdated(name, vpc.SubnetIDs[:2]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated App Block Builder"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated complex acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "enable_default_internet_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "prod"),
					resource.TestCheckResourceAttr(resourceName, "tags.Team", "platform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Environment", "prod"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Team", "platform"),
				),
			},
		},
	})
}

func TestAccAppBlockBuilder_noopPlan(t *testing.T) {
	testEnv := testhelpers.LoadAccTestEnv(t)
	vpc := testEnv.DefaultVPCInfo(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-builder")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBuilderBasicConfig(name, vpc.SubnetIDs[:2]),
			},
			{
				Config:             testAccAppBlockBuilderBasicConfig(name, vpc.SubnetIDs[:2]),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccAppBlockBuilderTagsConfig(name string, subnetIDs []string) string {
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
  }
}
`, name, testhelpers.HCLStringList(subnetIDs))
}

func TestAccAppBlockBuilder_tags(t *testing.T) {
	testEnv := testhelpers.LoadAccTestEnv(t)
	vpc := testEnv.DefaultVPCInfo(t)

	name := acctest.RandomWithPrefix("tf-acc-app-block-builder-tags")
	resourceName := "awsappstream_app_block_builder.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBuilderTagsConfig(name, vpc.SubnetIDs[:2]),
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
