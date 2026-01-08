// Copyright (c) St3ffn
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
	name := acctest.RandomWithPrefix("tf-acc-image-builder")
	resourceName := "awsappstream_image_builder.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
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
		},
	})
}

func testAccImageBuilderImageARNConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  image_arn     = "arn:aws:appstream:eu-central-1::image/AppStream-RockyLinux8-11-10-2025"
}
`, name)
}

func TestAccImageBuilder_imageARN(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-image-builder-arn")
	resourceName := "awsappstream_image_builder.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageBuilderImageARNConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_arn", "arn:aws:appstream:eu-central-1::image/AppStream-RockyLinux8-11-10-2025"),
					resource.TestCheckNoResourceAttr(resourceName, "image_name"),
				),
			},
		},
	})
}

func testAccImageBuilderRootVolumeConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  image_name    = "AppStream-RockyLinux8-11-10-2025"

  root_volume_config = {
    volume_size_in_gb = 250
  }
}
`, name)
}

func TestAccImageBuilder_rootVolumeConfig(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-image-builder-rootvol")
	resourceName := "awsappstream_image_builder.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageBuilderRootVolumeConfig(name),
				Check:  resource.TestCheckResourceAttr(resourceName, "root_volume_config.volume_size_in_gb", "250"),
			},
		},
	})
}

func testAccImageBuilderDescriptionTagsConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  image_name    = "AppStream-RockyLinux8-11-10-2025"

  description  = "test description"
  display_name = "Test Builder"

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`, name)
}

func TestAccImageBuilder_description(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-image-builder-desc")
	resourceName := "awsappstream_image_builder.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageBuilderDescriptionTagsConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Builder"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "terraform"),
				),
			},
		},
	})
}

func TestAccImageBuilder_noopPlan(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-image-builder-noop")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testAccImageBuilderBasicConfig(name)},
			{
				Config:             testAccImageBuilderBasicConfig(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
