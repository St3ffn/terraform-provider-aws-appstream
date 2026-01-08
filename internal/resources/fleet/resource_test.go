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

func testAccFleetBasicConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_fleet" "test" {
  name          = %q
  fleet_type   = "ON_DEMAND"
  instance_type = "stream.standard.small"

  image_name = "Amazon-AppStream2-Sample-Image-06-17-2024"

  compute_capacity = {
    desired_instances = 0
  }
}
`, name)
}

func TestAccFleet_basic(t *testing.T) {
	fleetName := acctest.RandomWithPrefix("tf-acc-fleet")

	resourceName := "awsappstream_fleet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetBasicConfig(fleetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fleetName),
					resource.TestCheckResourceAttr(resourceName, "fleet_type", "ON_DEMAND"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "stream.standard.small"),
					resource.TestCheckResourceAttr(resourceName, "compute_capacity.desired_instances", "0"),
					resource.TestCheckResourceAttr(resourceName, "image_name", "Amazon-AppStream2-Sample-Image-06-17-2024"),
					resource.TestCheckResourceAttr(resourceName, "image_arn", "arn:aws:appstream:eu-central-1::image/Amazon-AppStream2-Sample-Image-06-17-2024"),
					resource.TestCheckResourceAttr(resourceName, "enable_default_internet_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "max_user_duration_in_seconds", "57600"),
					resource.TestCheckResourceAttr(resourceName, "disconnect_timeout_in_seconds", "900"),
					resource.TestCheckResourceAttr(resourceName, "idle_disconnect_timeout_in_seconds", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "tags"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
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

func testAccFleetImageARNConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_fleet" "test" {
  name          = %q
  fleet_type    = "ON_DEMAND"
  instance_type = "stream.standard.small"

  image_arn = "arn:aws:appstream:eu-central-1::image/Amazon-AppStream2-Sample-Image-06-17-2024"

  compute_capacity = {
    desired_instances = 0
  }
}
`, name)
}

func TestAccFleet_imageARN(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-fleet-imgarn")
	resourceName := "awsappstream_fleet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetImageARNConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_arn", "arn:aws:appstream:eu-central-1::image/Amazon-AppStream2-Sample-Image-06-17-2024"),
					resource.TestCheckResourceAttr(resourceName, "image_name", "Amazon-AppStream2-Sample-Image-06-17-2024"),
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

func testAccFleetUpdateDescriptionTags(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_fleet" "test" {
  name          = %q
  fleet_type    = "ON_DEMAND"
  instance_type = "stream.standard.small"
  image_name    = "Amazon-AppStream2-Sample-Image-06-17-2024"

  description = "updated description"

  compute_capacity = {
    desired_instances = 0
  }

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`, name)
}

func TestAccFleet_updateDescription(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-fleet-update")
	resourceName := "awsappstream_fleet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testAccFleetBasicConfig(name)},
			{
				Config: testAccFleetUpdateDescriptionTags(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "terraform"),
				),
			},
		},
	})
}

func testAccFleetUpdateImageName(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_fleet" "test" {
  name          = %q
  fleet_type    = "ON_DEMAND"
  instance_type = "stream.standard.small"
  image_name    = "Amazon-AppStream2-Sample-Image-03-11-2023"

  description = "updated description"

  compute_capacity = {
    desired_instances = 0
  }
}
`, name)
}

func TestAccFleet_updateImageName(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-fleet-update")
	resourceName := "awsappstream_fleet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testAccFleetBasicConfig(name)},
			{
				Config: testAccFleetUpdateImageName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", "Amazon-AppStream2-Sample-Image-03-11-2023"),
					resource.TestCheckResourceAttr(resourceName, "image_arn", "arn:aws:appstream:eu-central-1::image/Amazon-AppStream2-Sample-Image-03-11-2023"),
				),
			},
		},
	})
}

func testAccFleetIdleTimeoutConfig(name string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_fleet" "test" {
  name          = %q
  fleet_type    = "ON_DEMAND"
  instance_type = "stream.standard.small"
  image_name    = "Amazon-AppStream2-Sample-Image-06-17-2024"

  idle_disconnect_timeout_in_seconds = 600

  compute_capacity = {
    desired_instances = 0
  }
}
`, name)
}

func TestAccFleet_idleTimeout(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-fleet-idle")
	resourceName := "awsappstream_fleet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetIdleTimeoutConfig(name),
				Check:  resource.TestCheckResourceAttr(resourceName, "idle_disconnect_timeout_in_seconds", "600"),
			},
		},
	})
}
