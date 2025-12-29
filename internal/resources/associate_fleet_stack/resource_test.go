// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_fleet_stack_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccAssociateFleetStackBasicConfig(fleetName, stackName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_stack" "test" {
  name = %q
}

resource "awsappstream_fleet" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  compute_capacity {
    desired_instances = 1
  }
}

resource "awsappstream_associate_fleet_stack" "test" {
  fleet_name = awsappstream_fleet.test.name
  stack_name = awsappstream_stack.test.name
}
`, stackName, fleetName)
}

func TestAccAssociateFleetStack_basic(t *testing.T) {
	stackName := acctest.RandomWithPrefix("tf-acc-stack")
	fleetName := acctest.RandomWithPrefix("tf-acc-fleet")

	resourceName := "awsappstream_associate_fleet_stack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateFleetStackBasicConfig(fleetName, stackName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "fleet_name", fleetName),
					resource.TestCheckResourceAttr(resourceName, "stack_name", stackName),
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
