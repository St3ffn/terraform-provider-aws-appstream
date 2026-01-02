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

func testAccUserWithDataSource(authenticationType, userName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_user" "test" {
  authentication_type = %q
  user_name           = %q
}

data "awsappstream_user" "test" {
  authentication_type = awsappstream_user.test.authentication_type
  user_name           = awsappstream_user.test.user_name
}
`, authenticationType, userName)
}

func TestAccUserDataSource_basic(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-acc-user") + "@example.com"
	resourceName := "data.awsappstream_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserWithDataSource("USERPOOL", userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "USERPOOL"),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}
