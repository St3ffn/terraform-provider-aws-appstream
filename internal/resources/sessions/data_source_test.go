// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package sessions_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func testAccSessionsDataSource_basic(stackName, fleetName string) string {
	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
data "awsappstream_sessions" "test" {
  stack_name = %q
  fleet_name = %q
}
`, stackName, fleetName)
}

func TestAccSoftwareAssociationsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionsDataSource_basic("fake-stack", "fake-fleet"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsappstream_sessions.test", "stack_name", "fake-stack"),
					resource.TestCheckResourceAttr("data.awsappstream_sessions.test", "fleet_name", "fake-fleet"),
					resource.TestCheckNoResourceAttr("data.awsappstream_sessions.test", "user_id"),
					resource.TestCheckNoResourceAttr("data.awsappstream_sessions.test", "authentication_type"),
					resource.TestCheckNoResourceAttr("data.awsappstream_sessions.test", "instance_id"),
					resource.TestCheckResourceAttr("data.awsappstream_sessions.test", "sessions.#", "0"),
				),
			},
		},
	})
}
