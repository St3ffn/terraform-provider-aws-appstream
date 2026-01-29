// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package export_image_tasks_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccExportImageTasksDataSource_basic() string {
	return testhelpers.TestAccProviderBasicConfig() + `
data "awsappstream_export_image_tasks" "test" {}
`
}

func TestAccExportImageTasksDataSource_basic(t *testing.T) {
	testhelpers.AccTestContextFromEnv(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportImageTasksDataSource_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_export_image_tasks.test", "export_image_tasks.#", "0"),
				),
			},
		},
	})
}

func testAccExportImageTasksDataSource_withFilter() string {
	return testhelpers.TestAccProviderBasicConfig() + `
data "awsappstream_export_image_tasks" "test" {
  filters = [{
    name   = "image-arn"
    values = ["test"]
  }]
}
`
}

func TestAccExportImageTasksDataSource_withFilter(t *testing.T) {
	testhelpers.AccTestContextFromEnv(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportImageTasksDataSource_withFilter(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_export_image_tasks.test", "export_image_tasks.#", "0"),
				),
			},
		},
	})
}
