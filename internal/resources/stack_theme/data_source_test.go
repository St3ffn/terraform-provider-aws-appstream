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

func testAccStackThemeWithDataSource(stackName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_stack_theme" "test" {
  stack_name   = awsappstream_stack.test.name
  title_text   = "Terraform DS Test"
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
      display_name    = "Docs"
      footer_link_url = "https://example.com/docs"
    }
  ]
}

data "awsappstream_stack_theme" "test" {
  stack_name = awsappstream_stack.test.name

  depends_on = [awsappstream_stack_theme.test]
}
`, stackName)
}

func TestAccStackThemeDataSource_basic(t *testing.T) {
	stackName := acctest.RandomWithPrefix("tf-acc-stack-ds")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackThemeWithDataSource(stackName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_stack_theme.test", "stack_name", stackName),
					resource.TestCheckResourceAttr("data.awsappstream_stack_theme.test", "title_text", "Terraform DS Test"),
					resource.TestCheckResourceAttr("data.awsappstream_stack_theme.test", "theme_styling", "BLUE"),
					resource.TestCheckResourceAttr("data.awsappstream_stack_theme.test", "footer_links.#", "1"),
					resource.TestCheckResourceAttr("data.awsappstream_stack_theme.test", "footer_links.0.display_name", "Docs"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack_theme.test", "state"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack_theme.test", "created_time"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack_theme.test", "theme_organization_logo_url"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack_theme.test", "theme_favicon_url"),
					resource.TestCheckResourceAttrSet("data.awsappstream_stack_theme.test", "id"),
				),
			},
		},
	})
}
