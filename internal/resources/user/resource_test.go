// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccUserResource(userName, extra string) string {
	return fmt.Sprintf(`
resource "awsappstream_user" "test" {
  authentication_type = "USERPOOL"
  user_name           = %q
%s
}
`, userName, extra)
}

func testAccUserBasicConfig(userName string) string {
	return testhelpers.TestAccProviderBasicConfig() +
		testAccUserResource(userName, "")
}

func TestAccUser_basic(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-acc-user") + "@example.com"
	resourceName := "awsappstream_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserBasicConfig(userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "USERPOOL"),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckNoResourceAttr(resourceName, "first_name"),
					resource.TestCheckNoResourceAttr(resourceName, "last_name"),
					resource.TestCheckNoResourceAttr(resourceName, "message_action"),
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

func testAccUserWithNamesConfig(userName string) string {
	return testhelpers.TestAccProviderBasicConfig() +
		testAccUserResource(userName, `
  first_name     = "Terraform"
  last_name      = "User"
  message_action = "SUPPRESS"
`)
}

func TestAccUser_namesAndMessageAction(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-acc-user") + "@example.com"
	resourceName := "awsappstream_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserWithNamesConfig(userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "Terraform"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "User"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccUserEnabledConfig(userName string, enabled bool) string {
	return testhelpers.TestAccProviderBasicConfig() +
		testAccUserResource(userName, fmt.Sprintf(`
  enabled = %t
`, enabled))
}

func TestAccUser_enabledToggle(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-acc-user") + "@example.com"
	resourceName := "awsappstream_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserEnabledConfig(userName, true),
				Check:  resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
			},
			{
				Config: testAccUserEnabledConfig(userName, false),
				Check:  resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
			},
			{
				Config: testAccUserEnabledConfig(userName, true),
				Check:  resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
			},
		},
	})
}

func TestAccUser_noOpPlan(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-acc-user") + "@example.com"

	config := testAccUserBasicConfig(userName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
