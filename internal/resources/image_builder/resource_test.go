// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccImageBuilderBasicConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  image_name    = "AppStream-RockyLinux8-11-10-2025"
}
`, name)
}

func TestAccImageBuilder_basic(t *testing.T) {
	testhelpers.LoadAccTestEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-image-builder")
	resourceName := "awsappstream_image_builder.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageBuilderBasicConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "stream.standard.small"),
					resource.TestCheckResourceAttr(resourceName, "image_name", "AppStream-RockyLinux8-11-10-2025"),
					resource.TestCheckResourceAttr(resourceName, "enable_default_internet_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "root_volume_config.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "image_arn"),
					resource.TestCheckNoResourceAttr(resourceName, "tags"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "platform"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_name"}, // image_name is not returned from aws
			},
			{
				Config:             testAccImageBuilderBasicConfig(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccImageBuilderComplexConfig(region, name string) string {
	return testhelpers.TestAccProviderTagsConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  image_arn     = "arn:aws:appstream:%s::image/AppStream-RockyLinux8-11-10-2025"

  description  = "test description"
  display_name = "Test Builder"

  root_volume_config = {
    volume_size_in_gb = 250
  }

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`, name, region)
}

func TestAccImageBuilder_complex(t *testing.T) {
	testEnv := testhelpers.LoadAccTestEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-image-builder-arn")
	resourceName := "awsappstream_image_builder.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageBuilderComplexConfig(testEnv.Region, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "stream.standard.small"),
					resource.TestCheckResourceAttr(resourceName, "image_arn", fmt.Sprintf("arn:aws:appstream:%s::image/AppStream-RockyLinux8-11-10-2025", testEnv.Region)),
					resource.TestCheckNoResourceAttr(resourceName, "image_name"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Builder"),
					resource.TestCheckResourceAttr(resourceName, "root_volume_config.volume_size_in_gb", "250"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.MANAGED_BY", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.BUILD_WITH", "love"),
				),
			},
		},
	})
}
