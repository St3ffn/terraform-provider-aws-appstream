// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccAssociateUserStackBasicConfig(stackName, userName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_user" "test" {
  user_name           = %q
  authentication_type = "USERPOOL"
}

resource "awsappstream_associate_user_stack" "test" {
  stack_name          = awsappstream_stack.test.name
  user_name           = awsappstream_user.test.user_name
  authentication_type = awsappstream_user.test.authentication_type
}
`, stackName, userName)
}

func TestAccAssociateUserStack_basic(t *testing.T) {
	stackName := acctest.RandomWithPrefix("tf-acc-stack")
	userName := acctest.RandomWithPrefix("tf-acc-user") + "@example.com"

	resourceName := "awsappstream_associate_user_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateUserStackBasicConfig(stackName, userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stack_name", stackName),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "USERPOOL"),
					resource.TestCheckResourceAttr(resourceName, "send_email_notification", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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

func testAccAssociateUserStackSendEmailConfig(stackName, userName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_user" "test" {
  user_name           = %q
  authentication_type = "USERPOOL"
}

resource "awsappstream_associate_user_stack" "test" {
  stack_name          = awsappstream_stack.test.name
  user_name           = awsappstream_user.test.user_name
  authentication_type = awsappstream_user.test.authentication_type

  send_email_notification = true
}
`, stackName, userName)
}

func TestAccAssociateUserStack_sendEmailNotification(t *testing.T) {
	stackName := acctest.RandomWithPrefix("tf-acc-stack")
	userName := acctest.RandomWithPrefix("tf-acc-user") + "@example.com"

	resourceName := "awsappstream_associate_user_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateUserStackSendEmailConfig(stackName, userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stack_name", stackName),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "USERPOOL"),
					resource.TestCheckResourceAttr(resourceName, "send_email_notification", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
