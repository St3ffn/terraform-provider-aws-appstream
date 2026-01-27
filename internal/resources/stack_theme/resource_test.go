// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_theme_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccStackThemeBasicConfig(stackName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_stack_theme" "test" {
  stack_name  = awsappstream_stack.test.name
  title_text = "Terraform Acceptance Test"
  theme_styling = "BLUE"

  organization_logo_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  favicon_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }
}
`, stackName)
}

func TestAccStackTheme_basic(t *testing.T) {
	stackName := acctest.RandomWithPrefix("tf-acc-stack")

	resourceName := "awsappstream_stack_theme.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackThemeBasicConfig(stackName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stack_name", stackName),
					resource.TestCheckResourceAttr(resourceName, "title_text", "Terraform Acceptance Test"),
					resource.TestCheckResourceAttr(resourceName, "theme_styling", "BLUE"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "theme_organization_logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "theme_favicon_url"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"organization_logo_s3_location",
					"favicon_s3_location",
					"theme_organization_logo_url",
					"theme_favicon_url",
				},
			},
			{
				Config:             testAccStackThemeBasicConfig(stackName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccStackThemeComplexConfig(stackName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_stack_theme" "test" {
  stack_name   = awsappstream_stack.test.name
  title_text  = "Terraform Complex Test"
  theme_styling = "BLUE"

  organization_logo_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  favicon_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  footer_links = [
    {
      display_name     = "Documentation"
      footer_link_url  = "https://example.com/docs"
    },
    {
      display_name     = "Support"
      footer_link_url  = "https://example.com/support"
    }
  ]
}
`, stackName)
}

func testAccStackThemeComplexConfigUpdated(stackName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_stack_theme" "test" {
  stack_name   = awsappstream_stack.test.name
  title_text  = "Terraform Complex Test Updated"
  theme_styling = "BLUE"

  organization_logo_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  favicon_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  footer_links = [
    {
      display_name     = "Docs"
      footer_link_url  = "https://example.com/docs-v2"
    }
  ]
}
`, stackName)
}

func testAccStackThemeComplexConfigNoFooter(stackName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_stack_theme" "test" {
  stack_name   = awsappstream_stack.test.name
  title_text  = "Terraform Complex Test Updated"
  theme_styling = "BLUE"

  organization_logo_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }

  favicon_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
    s3_key    = "application_icon.png"
  }
}
`, stackName)
}

func TestAccStackTheme_complex(t *testing.T) {
	stackName := acctest.RandomWithPrefix("tf-acc-stack")

	resourceName := "awsappstream_stack_theme.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackThemeComplexConfig(stackName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "footer_links.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "footer_links.0.display_name", "Documentation"),
					resource.TestCheckResourceAttr(resourceName, "footer_links.1.display_name", "Support"),
				),
			},
			{
				Config: testAccStackThemeComplexConfigUpdated(stackName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title_text", "Terraform Complex Test Updated"),
					resource.TestCheckResourceAttr(resourceName, "footer_links.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "footer_links.0.display_name", "Docs"),
				),
			},
			{
				Config: testAccStackThemeComplexConfigNoFooter(stackName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "footer_links")),
			},
		},
	})
}
