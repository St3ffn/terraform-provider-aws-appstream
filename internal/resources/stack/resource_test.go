// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccStackResource(name string, extra string) string {
	return fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
%s
}
`, name, extra)
}

func testAccStackBasicConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() +
		testAccStackResource(name, "")
}

func TestAccStack_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-basic")
	resourceName := "awsappstream_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackBasicConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckNoResourceAttr(resourceName, "tags"),
					resource.TestCheckNoResourceAttr(resourceName, "user_settings"),
					resource.TestCheckNoResourceAttr(resourceName, "application_settings"),
					resource.TestCheckNoResourceAttr(resourceName, "streaming_experience_settings"),
					resource.TestCheckNoResourceAttr(resourceName, "access_endpoints"),
					resource.TestCheckNoResourceAttr(resourceName, "storage_connectors"),
					resource.TestCheckNoResourceAttr(resourceName, "embed_host_domains"),
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

func testAccStackTagsConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() +
		testAccStackResource(name, `
  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
`)
}

func TestAccStack_tags(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-tags")
	resourceName := "awsappstream_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackTagsConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "terraform"),
				),
			},
		},
	})
}

func testAccStackDescriptionConfig(name, description string) string {
	return testhelpers.TestAccProviderBasicConfig() +
		testAccStackResource(name, fmt.Sprintf(`
  description = %q
`, description))
}

func TestAccStack_updateDescription(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-update")
	resourceName := "awsappstream_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testAccStackBasicConfig(name)},
			{
				Config: testAccStackDescriptionConfig(name, "updated"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "updated"),
			},
		},
	})
}

func testAccStackComplexConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q

  redirect_url = "https://example.com/logout"

  streaming_experience_settings = {
    preferred_protocol = "TCP"
  }

  application_settings = {
    enabled = false
  }

  user_settings = [
    {
      action     = "CLIPBOARD_COPY_FROM_LOCAL_DEVICE"
      permission = "ENABLED"
    }
  ]
}
`, name)
}

func testAccStackComplexConfigUpdated(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q

  redirect_url = "https://example.com/updated"

  streaming_experience_settings = {
    preferred_protocol = "UDP"
  }

  application_settings = {
    enabled        = true
    settings_group = "test-group"
  }

  user_settings = [
    {
      action     = "FILE_UPLOAD"
      permission = "DISABLED"
    }
  ]
}
`, name)
}

func TestAccStack_complex(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-complex")
	resourceName := "awsappstream_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackComplexConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://example.com/logout"),
					resource.TestCheckResourceAttr(resourceName, "streaming_experience_settings.preferred_protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "application_settings.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "user_settings.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"user_settings.*",
						map[string]string{
							"action":     "CLIPBOARD_COPY_FROM_LOCAL_DEVICE",
							"permission": "ENABLED",
						},
					),
					resource.TestCheckNoResourceAttr(resourceName, "storage_connectors"),
					resource.TestCheckNoResourceAttr(resourceName, "access_endpoints"),
					resource.TestCheckNoResourceAttr(resourceName, "embed_host_domains"),
				),
			},
			{
				Config: testAccStackComplexConfigUpdated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://example.com/updated"),
					resource.TestCheckResourceAttr(resourceName, "streaming_experience_settings.preferred_protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceName, "application_settings.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "application_settings.settings_group", "test-group"),
					resource.TestCheckResourceAttr(resourceName, "user_settings.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"user_settings.*",
						map[string]string{
							"action":     "FILE_UPLOAD",
							"permission": "DISABLED",
						},
					),
				),
			},
		},
	})
}

func testAccStackStorageConnectorConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q

  storage_connectors = [
    {
      connector_type = "HOMEFOLDERS"
    }
  ]
}
`, name)
}

func TestAccStack_storageConnector(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-storage")
	resourceName := "awsappstream_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackStorageConnectorConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "storage_connectors.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"storage_connectors.*",
						map[string]string{
							"connector_type": "HOMEFOLDERS",
						},
					),
				),
			},
			{
				Config:             testAccStackStorageConnectorConfig(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccStackUserSettingsMaxLength(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q

  user_settings = [
    {
      action         = "CLIPBOARD_COPY_TO_LOCAL_DEVICE"
      permission     = "ENABLED"
      maximum_length = 1024
    }
  ]
}
`, name)
}

func TestAccStack_userSettingsMaximumLength(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-maxlen")
	resourceName := "awsappstream_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackUserSettingsMaxLength(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_settings.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"user_settings.*",
						map[string]string{
							"action":         "CLIPBOARD_COPY_TO_LOCAL_DEVICE",
							"permission":     "ENABLED",
							"maximum_length": "1024",
						},
					),
				),
			},
		},
	})
}

func testAccStackEmbedDomainsConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q

  embed_host_domains = ["example.com"]
}
`, name)
}

func TestAccStack_embedHostDomains(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-stack-embed")
	resourceName := "awsappstream_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackEmbedDomainsConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "embed_host_domains.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName,
						"embed_host_domains.*",
						"example.com",
					),
				),
			},
		},
	})
}
