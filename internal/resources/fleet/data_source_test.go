// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccFleetWithDataSource(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_fleet" "test" {
  name          = %q
  fleet_type    = "ON_DEMAND"
  instance_type = "stream.standard.small"

  image_name = "Amazon-AppStream2-Sample-Image-06-17-2024"

  compute_capacity = {
    desired_instances = 0
  }
}

data "awsappstream_fleet" "test" {
  name = awsappstream_fleet.test.name
}
`, name)
}

func TestAccFleetDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-fleet-ds-basic")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetWithDataSource(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.awsappstream_fleet.test", "name", name,
					),
					resource.TestCheckResourceAttr(
						"data.awsappstream_fleet.test", "fleet_type", "ON_DEMAND",
					),
					resource.TestCheckResourceAttr(
						"data.awsappstream_fleet.test", "enable_default_internet_access", "false"),
					resource.TestCheckResourceAttr(
						"data.awsappstream_fleet.test", "max_user_duration_in_seconds", "57600"),
					resource.TestCheckResourceAttr(
						"data.awsappstream_fleet.test", "disconnect_timeout_in_seconds", "900"),
					resource.TestCheckResourceAttr(""+
						"data.awsappstream_fleet.test", "idle_disconnect_timeout_in_seconds", "0"),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_fleet.test", "arn",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_fleet.test", "created_time",
					),
					resource.TestCheckResourceAttrSet(
						"data.awsappstream_fleet.test", "state",
					),
				),
			},
		},
	})
}

func TestAccFleetDataSource_computed(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-fleet-ds-computed")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetWithDataSource(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.awsappstream_fleet.test",
						"compute_capacity.desired_instances",
						"0",
					),
					resource.TestCheckNoResourceAttr(
						"data.awsappstream_fleet.test",
						"vpc_config",
					),
					resource.TestCheckNoResourceAttr(
						"data.awsappstream_fleet.test",
						"domain_join_info",
					),
				),
			},
		},
	})
}
