// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccUsageReportSubscriptionBasicConfig() string {
	return testhelpers.TestAccProviderBasicConfig() + `
resource "awsappstream_usage_report_subscription" "test" {}
`
}

func TestAccUsageReportSubscription_basic(t *testing.T) {
	testhelpers.AccTestContextFromEnv(t)

	resourceName := "awsappstream_usage_report_subscription.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUsageReportSubscriptionBasicConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "usage-report-subscription"),
					resource.TestCheckResourceAttrSet(resourceName, "s3_bucket_name"),
					resource.TestCheckResourceAttr(resourceName, "schedule", "DAILY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "ignored",
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "usage-report-subscription"),
				),
			},
			{
				Config:             testAccUsageReportSubscriptionBasicConfig(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
