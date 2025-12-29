// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccEntitlementWithDataSource(name, stackName, description string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_entitlement" "test" {
  stack_name     = awsappstream_stack.test.name
  name           = %q
  description    = %q
  app_visibility = "ALL"

  attributes = [
    {
      name  = "roles"
      value = "test"
    }
  ]
}

data "awsappstream_entitlement" "test" {
  stack_name = awsappstream_stack.test.name
  name       = awsappstream_entitlement.test.name
}
`, stackName, name, description)
}

func TestAccEntitlementDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-entitlement-ds")
	stackName := acctest.RandomWithPrefix("tf-acc-entitlement-stack-ds")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEntitlementWithDataSource(name, stackName, "Acceptance test entitlement"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_entitlement.test", "name", name),
					resource.TestCheckResourceAttr("data.awsappstream_entitlement.test", "stack_name", stackName),
					resource.TestCheckResourceAttr("data.awsappstream_entitlement.test", "app_visibility", "ALL"),
					resource.TestCheckResourceAttr("data.awsappstream_entitlement.test", "attributes.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.awsappstream_entitlement.test",
						"attributes.*",
						map[string]string{
							"name":  "roles",
							"value": "test",
						},
					),
					resource.TestCheckResourceAttrSet("data.awsappstream_entitlement.test", "created_time"),
				),
			},
		},
	})
}

func TestAccEntitlementDataSource_withDescription(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-entitlement-ds-desc")
	stackName := acctest.RandomWithPrefix("tf-acc-entitlement-stack-ds-desc")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEntitlementWithDataSource(
					name,
					stackName,
					"entitlement description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.awsappstream_entitlement.test",
						"description",
						"entitlement description",
					),
				),
			},
		},
	})
}
