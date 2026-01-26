// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package software_associations_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/testhelpers"
)

func getAWSAccountInfo(t *testing.T) (accountID, region string) {
	if !testhelpers.IsAccTest() {
		t.Skip("Skipping acceptance test unless TF_ACC is set")
	}

	t.Helper()

	cfg, err := config.LoadDefaultConfig(t.Context())
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	out, err := stsClient.GetCallerIdentity(t.Context(), &sts.GetCallerIdentityInput{})
	if err != nil {
		t.Fatalf("failed to get caller identity: %v", err)
	}

	if out.Account == nil {
		t.Fatalf("STS GetCallerIdentity returned nil account ID")
	}

	if cfg.Region == "" {
		t.Fatalf("AWS region not set")
	}

	return *out.Account, cfg.Region
}

func testAccSoftwareAssociationsDataSource_basic(t *testing.T) string {
	accountID, region := getAWSAccountInfo(t)

	return testhelpers.TestAccProviderBasicConfig() + fmt.Sprintf(`
data "awsappstream_software_associations" "test" {
  associated_resource = "arn:aws:appstream:%s:%s:image/fake-image-for-test"
}
`, region, accountID)
}

func TestAccSoftwareAssociationsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testhelpers.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testhelpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSoftwareAssociationsDataSource_basic(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.awsappstream_software_associations.test", "associated_resource"),
					resource.TestCheckResourceAttr("data.awsappstream_software_associations.test", "software_associations.#", "0"),
				),
			},
		},
	})
}
