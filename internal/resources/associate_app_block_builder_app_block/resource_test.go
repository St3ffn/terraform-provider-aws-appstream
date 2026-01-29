// Copyright St3ffn 2025, 2026
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
	bucketName string,
) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
resource "awsappstream_app_block" "test" {
  name = %[2]q

  packaging_type = "APPSTREAM2"

  source_s3_location = {
    s3_bucket = %[4]q
  }
}

resource "awsappstream_app_block_builder" "test" {
  name          = %[1]q
  instance_type = "stream.standard.small"
  platform      = "WINDOWS_SERVER_2019"

  vpc_config = {
    subnet_ids = [%[3]s]
  }
}

resource "awsappstream_associate_app_block_builder_app_block" "test" {
  app_block_builder_name = awsappstream_app_block_builder.test.name
  app_block_arn          = awsappstream_app_block.test.arn
}
`,
		builderName,
		appBlockName,
		testhelpers.HCLStringList(subnetIDs),
		bucketName,
	)
}

func TestAccAppBlockBuilderAssociation_basic(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)
	vpc := testCtx.DefaultVPCInfo(t)

	builderName := acctest.RandomWithPrefix("tf-acc-app-block-builder")
	appBlockName := acctest.RandomWithPrefix("tf-acc-app-block")

	resourceName := "awsappstream_associate_app_block_builder_app_block.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBlockBuilderAssociationBasicConfig(builderName, appBlockName, vpc.SubnetIDs[:2], testCtx.BucketName),
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
				Config:             testAccAppBlockBuilderAssociationBasicConfig(builderName, appBlockName, vpc.SubnetIDs[:2], testCtx.BucketName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
