// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccUsageReportSubscriptionDataSourceBasic() string {
	return testhelpers.TestAccProviderBasicConfig() + `
resource "awsappstream_usage_report_subscription" "test" {}

data "awsappstream_usage_report_subscription" "test" {
  depends_on = [awsappstream_usage_report_subscription.test]
}
`
}

func TestAccUsageReportSubscriptionDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUsageReportSubscriptionDataSourceBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.awsappstream_usage_report_subscription.test", "s3_bucket_name"),
					resource.TestCheckResourceAttr("data.awsappstream_usage_report_subscription.test", "schedule", "DAILY"),
				),
			},
		},
	})
}
