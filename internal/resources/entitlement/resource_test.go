// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccEntitlementBasicConfig(name, stackName, description string) string {
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
`, stackName, name, description)
}

/*
 * tests
 */

func TestAccEntitlement_basic(t *testing.T) {
	testhelpers.AccTestContextFromEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-entitlement")
	stackName := acctest.RandomWithPrefix("tf-acc-entitlement-stack")
	resourceName := "awsappstream_entitlement.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEntitlementBasicConfig(name, stackName, "Acceptance test entitlement"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "stack_name", stackName),
					resource.TestCheckResourceAttr(resourceName, "app_visibility", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "attributes.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"attributes.*",
						map[string]string{
							"name":  "roles",
							"value": "test",
						},
					),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccEntitlementBasicConfig(name, stackName, "Acceptance test entitlement"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccEntitlement_updateDescription(t *testing.T) {
	testhelpers.AccTestContextFromEnv(t)

	name := acctest.RandomWithPrefix("tf-acc-entitlement-update")
	stackName := acctest.RandomWithPrefix("tf-acc-entitlement-stack-update")
	resourceName := "awsappstream_entitlement.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEntitlementBasicConfig(name, stackName, "initial description"),
			},
			{
				Config: testAccEntitlementBasicConfig(name, stackName, "updated description"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
			},
		},
	})
}
