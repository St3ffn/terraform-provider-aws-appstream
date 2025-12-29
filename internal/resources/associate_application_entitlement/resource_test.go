// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_entitlement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccAssociateApplicationEntitlementBasicConfig(
	stackName, entitlementName, applicationName string,
) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_entitlement" "test" {
  stack_name     = awsappstream_stack.test.name
  name           = %q
  app_visibility = "ASSOCIATED"

  attributes = [{
    name  = "title"
    value = "test"
  }]
}

resource "awsappstream_associate_application_entitlement" "test" {
  stack_name             = awsappstream_stack.test.name
  entitlement_name       = awsappstream_entitlement.test.name
  application_identifier = %q
}
`, stackName, entitlementName, applicationName)
}

func TestAccAssociateApplicationEntitlement_basic(t *testing.T) {
	stackName := acctest.RandomWithPrefix("tf-acc-stack")
	entitlementName := acctest.RandomWithPrefix("tf-acc-entitlement")
	applicationName := acctest.RandomWithPrefix("tf-acc-app")

	resourceName := "awsappstream_associate_application_entitlement.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateApplicationEntitlementBasicConfig(stackName, entitlementName, applicationName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stack_name", stackName),
					resource.TestCheckResourceAttr(resourceName, "entitlement_name", entitlementName),
					resource.TestCheckResourceAttrSet(resourceName, "application_identifier"),
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
