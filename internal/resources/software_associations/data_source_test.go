// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package software_associations_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccSoftwareAssociationsDataSource_basic(region, accountID string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
data "awsappstream_software_associations" "test" {
  associated_resource = "arn:aws:appstream:%s:%s:image/fake-image-for-test"
}
`, region, accountID)
}

func TestAccSoftwareAssociationsDataSource_basic(t *testing.T) {
	testCtx := testhelpers.AccTestContextFromEnv(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSoftwareAssociationsDataSource_basic(testCtx.Region, testCtx.AccountID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.awsappstream_software_associations.test", "associated_resource"),
					resource.TestCheckResourceAttr("data.awsappstream_software_associations.test", "software_associations.#", "0"),
				),
			},
		},
	})
}
