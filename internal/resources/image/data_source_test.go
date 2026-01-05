// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

const testAccPublicImageName = "Amazon-AppStream2-Sample-Image-06-17-2024"
const testAccPublicImageARN = "arn:aws:appstream:eu-central-1::image/" + testAccPublicImageName

func testAccImageDataSourceByARN() string {
	return testhelpers.TestAccProviderBasicConfig() + `
data "awsappstream_image" "test" {
  arn = "` + testAccPublicImageARN + `"
}
`
}

func TestAccImageDataSource_basicByARN(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceByARN(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"arn",
						testAccPublicImageARN,
					),
					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"name",
						testAccPublicImageName,
					),
					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"visibility",
						"PUBLIC",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"created_time",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"platform",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"state",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"image_type",
					),
				),
			},
		},
	})
}

func testAccImageDataSourceByName() string {
	return testhelpers.TestAccProviderBasicConfig() + `
data "awsappstream_image" "test" {
  name = "` + testAccPublicImageName + `"
}
`
}

func TestAccImageDataSource_basicByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceByName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"arn",
						testAccPublicImageARN,
					),
					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"name",
						testAccPublicImageName,
					),
					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"visibility",
						"PUBLIC",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"created_time",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"platform",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"state",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"image_type",
					),
				),
			},
		},
	})
}
func testAccImageDataSourceByRegexMostRecent() string {
	return testhelpers.TestAccProviderBasicConfig() + `
data "awsappstream_image" "test" {
  name_regex = "^Amazon-AppStream2-Sample-Image-"
  most_recent = true
}
`
}

func TestAccImageDataSource_byRegexMostRecent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceByRegexMostRecent(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"arn",
						testAccPublicImageARN,
					),
					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"name",
						testAccPublicImageName,
					),

					resource.TestCheckResourceAttr(
						"data.awsappstream_image.test",
						"visibility",
						"PUBLIC",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"created_time",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"platform",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"state",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_image.test",
						"image_type",
					),
				),
			},
		},
	})
}

func testAccImageDataSourceByRegexMultiple() string {
	return testhelpers.TestAccProviderBasicConfig() + `
data "awsappstream_image" "test" {
  name_regex = "^Amazon-AppStream2-Sample-Image-"
}
`
}

func TestAccImageDataSource_byRegexMultipleFails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccImageDataSourceByRegexMultiple(),
				ExpectError: regexp.MustCompile(`multiple images matched the selection criteria`),
			},
		},
	})
}
