// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccAppBlockBuilderAssociationBasicConfig(
	builderName string,
	appBlockName string,
	subnetIDs []string,
) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %q

  packaging_type = "APPSTREAM2"

  source_s3_location = {
    s3_bucket = "appstream-acc-test-bucket"
  }
}

resource "awsappstream_app_block_builder" "test" {
  name          = %q
  instance_type = "stream.standard.small"
  platform      = "WINDOWS_SERVER_2019"

  vpc_config = {
    subnet_ids = [%s]
  }
}

resource "awsappstream_associate_app_block_builder_app_block" "test" {
  app_block_builder_name = awsappstream_app_block_builder.test.name
  app_block_arn          = awsappstream_app_block.test.arn
}
`,
		appBlockName,
		builderName,
		testhelpers.HCLStringList(subnetIDs),
	)
}

func TestAccAppBlockBuilderAssociation_basic(t *testing.T) {
	vpcInfo, err := testhelpers.GetDefaultVPCInfo(t)
	if err != nil {
		t.Fatalf("failed to get default VPC info: %v", err)
	}

	builderName := acctest.RandomWithPrefix("tf-acc-app-block-builder")
	appBlockName := acctest.RandomWithPrefix("tf-acc-app-block")

	resourceName := "awsappstream_associate_app_block_builder_app_block.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBuilderAssociationBasicConfig(builderName, appBlockName, vpcInfo.SubnetIDs[:2]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app_block_builder_name", builderName),
					resource.TestCheckResourceAttrSet(resourceName, "app_block_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccAppBlockBuilderAssociationBasicConfig(builderName, appBlockName, vpcInfo.SubnetIDs[:2]),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
