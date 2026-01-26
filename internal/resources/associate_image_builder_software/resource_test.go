// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccAssociateImageBuilderSoftwareBasicConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.medium"
  image_name    = "AppStream-WinServer2025-12-18-2025"
}

resource "awsappstream_associate_image_builder_software" "test" {
  image_builder_arn = awsappstream_image_builder.test.arn
  software_names = [
    "Microsoft_Office_2024_LTSC_Standard_64Bit"
  ]
}
`, name)
}

func TestAccAssociateImageBuilderSoftware_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-imgb")

	resourceName := "awsappstream_associate_image_builder_software.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateImageBuilderSoftwareBasicConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "software_names.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName,
						"software_names.0",
						"Microsoft_Office_2024_LTSC_Standard_64Bit",
					),
					resource.TestCheckResourceAttr(resourceName, "deploy", "false"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// import only establishes the association identity (image_builder_arn).
				// software_names and deploy represent user intent and are not returned by AWS.
				// associations is informational-only and cannot be reconstructed on import
				// without declared software_names, so these attributes are intentionally ignored.
				ImportStateVerifyIgnore: []string{
					"software_names",
					"associations",
					"deploy",
				},
			},
			{
				Config:             testAccAssociateImageBuilderSoftwareBasicConfig(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccAssociateImageBuilderSoftwareUpdateConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.medium"
  image_name    = "AppStream-WinServer2025-12-18-2025"
}

resource "awsappstream_associate_image_builder_software" "test" {
  image_builder_arn = awsappstream_image_builder.test.arn
  software_names = [
    "Microsoft_Office_2024_LTSC_Standard_64Bit",
    "Microsoft_Project_2024_Standard_64Bit"
  ]
}
`, name)
}

func TestAccAssociateImageBuilderSoftware_update(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-imgb-update")
	resourceName := "awsappstream_associate_image_builder_software.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateImageBuilderSoftwareBasicConfig(name),
			},
			{
				Config: testAccAssociateImageBuilderSoftwareUpdateConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "software_names.#", "2"),
				),
			},
		},
	})
}

func testAccAssociateImageBuilderSoftwareDeployConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_image_builder" "test" {
  name          = %q
  instance_type = "stream.standard.medium"
  image_name    = "AppStream-WinServer2025-12-18-2025"
}

resource "awsappstream_associate_image_builder_software" "test" {
  image_builder_arn = awsappstream_image_builder.test.arn
  software_names = [
    "Microsoft_Office_2024_LTSC_Standard_64Bit"
  ]
  deploy = true
}
`, name)
}

func TestAccAssociateImageBuilderSoftware_deploy(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-imgb-deploy")
	resourceName := "awsappstream_associate_image_builder_software.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateImageBuilderSoftwareDeployConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "deploy", "true"),
				),
			},
		},
	})
}
