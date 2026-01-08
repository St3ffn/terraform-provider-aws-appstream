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

func testAccImageBuilderWithDataSource(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  image_name    = "AppStream-RockyLinux8-11-10-2025"
}

data "awsappstream_image_builder" "test" {
  name = awsappstream_image_builder.test.name
}
`, name)
}

func TestAccImageBuilderDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-image-builder-ds-basic")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageBuilderWithDataSource(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_image_builder.test", "name", name),
					resource.TestCheckNoResourceAttr("data.awsappstream_image_builder.test", "tags"),
					resource.TestCheckResourceAttrSet("data.awsappstream_image_builder.test", "arn"),
					resource.TestCheckResourceAttrSet("data.awsappstream_image_builder.test", "created_time"),
					resource.TestCheckResourceAttrSet("data.awsappstream_image_builder.test", "platform"),
					resource.TestCheckResourceAttrSet("data.awsappstream_image_builder.test", "state"),
					resource.TestCheckResourceAttr("data.awsappstream_image_builder.test", "root_volume_config.%", "1"),
				),
			},
		},
	})
}

func testAccImageBuilderWithDataSourceArnDescriptionTags(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  image_arn     = "arn:aws:appstream:eu-central-1::image/AppStream-RockyLinux8-11-10-2025"

  description  = "test description"
  display_name = "Test Builder"

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}

data "awsappstream_image_builder" "test" {
  name = awsappstream_image_builder.test.name
}
`, name)
}

func TestAccImageBuilderDataSource_arn_and_description(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-image-builder-ds-desc")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageBuilderWithDataSourceArnDescriptionTags(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_image_builder.test", "image_arn", "arn:aws:appstream:eu-central-1::image/AppStream-RockyLinux8-11-10-2025"),
					resource.TestCheckNoResourceAttr("data.awsappstream_image_builder.test", "image_name"),
					resource.TestCheckResourceAttr("data.awsappstream_image_builder.test", "description", "test description"),
					resource.TestCheckResourceAttr("data.awsappstream_image_builder.test", "display_name", "Test Builder"),
					resource.TestCheckResourceAttr("data.awsappstream_image_builder.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("data.awsappstream_image_builder.test", "tags.Owner", "terraform"),
				),
			},
		},
	})
}
