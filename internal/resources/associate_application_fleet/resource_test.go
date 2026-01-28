// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_fleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccAssociateApplicationFleetBasicConfig(
	appBlockName, applicationName, fleetName string, subnetIDs []string, bucketName string,
) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %[1]q

  packaging_type = "CUSTOM"

  source_s3_location = {
    s3_bucket = %[5]q
    s3_key    = "appblock.vhd"
  }

  setup_script_details = {
    script_s3_location = {
      s3_bucket = %[5]q
      s3_key    = "app_block_setup.sh"
    }

    executable_path    = "/bin/bash"
    timeout_in_seconds = 60
  }
}

resource "awsappstream_application" "test" {
  name = %[2]q

  icon_s3_location = {
    s3_bucket = %[5]q
    s3_key    = "application_icon.png"
  }

  launch_path = "/app/app"

  platforms = ["AMAZON_LINUX2"]
  instance_families = ["GENERAL_PURPOSE"]

  app_block_arn = awsappstream_app_block.test.arn
}

resource "awsappstream_fleet" "test" {
  name          = %[3]q
  fleet_type    = "ELASTIC"
  instance_type = "stream.standard.small"

  vpc_config = {
    subnet_ids = [%[4]s]
  }

  platform = "AMAZON_LINUX2"

  max_concurrent_sessions = 1
}

resource "awsappstream_associate_application_fleet" "test" {
  fleet_name      = awsappstream_fleet.test.name
  application_arn = awsappstream_application.test.arn
}
`, appBlockName, applicationName, fleetName, testhelpers.HCLStringList(subnetIDs), bucketName)
}

func TestAccAssociateApplicationFleet_basic(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)
	vpc := testCtx.DefaultVPCInfo(t)

	appBlockName := acctest.RandomWithPrefix("tf-acc-app-block")
	applicationName := acctest.RandomWithPrefix("tf-acc-application")
	fleetName := acctest.RandomWithPrefix("tf-acc-fleet")

	resourceName := "awsappstream_associate_application_fleet.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateApplicationFleetBasicConfig(appBlockName, applicationName, fleetName, vpc.SubnetIDs[:2], testCtx.BucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "fleet_name", fleetName),
					resource.TestCheckResourceAttrSet(resourceName, "application_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccAssociateApplicationFleetBasicConfig(appBlockName, applicationName, fleetName, vpc.SubnetIDs[:2], testCtx.BucketName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
