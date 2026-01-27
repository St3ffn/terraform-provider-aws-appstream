// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccDirectoryConfigBasic(name, ou, accountName, password string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_directory_config" "test" {
  directory_name = %q

  organizational_unit_distinguished_names = [%q]

  service_account_credentials = {
    account_name     = %q
    account_password = %q
  }
}
`, name, ou, accountName, password)
}

func TestAccDirectoryConfig_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-directory-config") + ".example.com"
	resourceName := "awsappstream_directory_config.test"

	ou := fmt.Sprintf("OU=Test,DC=%s", strings.ReplaceAll(name, ".", ",DC="))
	accountName := fmt.Sprintf("%s\\%s", name, "username")
	password := "password"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectoryConfigBasic(name, ou, accountName, password),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "directory_name", name),
					resource.TestCheckResourceAttr(resourceName, "organizational_unit_distinguished_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "organizational_unit_distinguished_names.0", ou),
					resource.TestCheckResourceAttr(resourceName, "service_account_credentials.account_name", accountName),
					resource.TestCheckResourceAttr(resourceName, "service_account_credentials.account_password", password),
					resource.TestCheckNoResourceAttr(resourceName, "service_account_credentials.certificate_based_auth_properties"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"service_account_credentials",
				},
			},
			{
				Config:             testAccDirectoryConfigBasic(name, ou, accountName, password),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccDirectoryConfig_complex(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-directory-config") + ".example.com"
	resourceName := "awsappstream_directory_config.test"

	ou1 := fmt.Sprintf("OU=Test1,DC=%s", strings.ReplaceAll(name, ".", ",DC="))
	ou2 := fmt.Sprintf("OU=Test2,DC=%s", strings.ReplaceAll(name, ".", ",DC="))

	accountName1 := fmt.Sprintf("%s\\user1", name)
	password1 := "password1"

	accountName2 := fmt.Sprintf("%s\\user2", name)
	password2 := "password2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectoryConfigBasic(name, ou1, accountName1, password1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "directory_name", name),
					resource.TestCheckResourceAttr(resourceName, "organizational_unit_distinguished_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "organizational_unit_distinguished_names.0", ou1),
					resource.TestCheckResourceAttr(resourceName, "service_account_credentials.account_name", accountName1),
					resource.TestCheckResourceAttr(resourceName, "service_account_credentials.account_password", password1),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
				),
			},
			{
				Config: testAccDirectoryConfigBasic(name, ou2, accountName2, password2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "directory_name", name),
					resource.TestCheckResourceAttr(resourceName, "organizational_unit_distinguished_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "organizational_unit_distinguished_names.0", ou2),
					resource.TestCheckResourceAttr(resourceName, "service_account_credentials.account_name", accountName2),
					resource.TestCheckResourceAttr(resourceName, "service_account_credentials.account_password", password2),
				),
			},
			{
				Config:             testAccDirectoryConfigBasic(name, ou2, accountName2, password2),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
