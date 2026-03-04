// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// BucketPrefix is the prefix used for S3 buckets created and referenced by
// acceptance tests. The final bucket name is suffixed with AWS account ID
// and region to avoid collisions across environments.
const BucketPrefix = "appstream-acc-test-bucket"

// AccTestEnv contains common, environment-derived information used by
// Terraform acceptance tests.
//
// It centralizes AWS-related context such as region, account ID, and shared
// resource names to avoid repeated environment lookups across tests.
type AccTestEnv struct {
	// Region is the AWS region resolved from the environment and SDK config.
	Region string
	// AccountID is the AWS account ID resolved via STS GetCallerIdentity.
	AccountID string
	// BucketName is the name of the shared S3 bucket used by acceptance tests.
	// It is derived from BucketPrefix, AccountID, and Region.
	BucketName string
}

// LoadAccTestEnv initializes and returns an AccTestEnv based on the
// current process environment and AWS SDK configuration.
//
// The function:
//   - Skips the test if TF_ACC is not set
//   - Verifies AWS credentials and region environment variables are present
//   - Loads the default AWS SDK configuration
//   - Resolves the AWS account ID via STS
//
// This helper should be called once per acceptance test and reused for
// all test configuration and assertions.
func LoadAccTestEnv(t *testing.T) *AccTestEnv {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("skipping acceptance test unless TF_ACC is set")
	}

	if os.Getenv("AWS_PROFILE") == "" &&
		os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Fatal("AWS credentials not set")
	}

	if os.Getenv("AWS_REGION") == "" &&
		os.Getenv("AWS_DEFAULT_REGION") == "" {
		t.Fatal("AWS region not set")
	}

	cfg, err := config.LoadDefaultConfig(t.Context())
	if err != nil {
		t.Fatalf("failed to load aws config: %v", err)
	}

	if cfg.Region == "" {
		t.Fatal("aws region not set")
	}

	stsClient := sts.NewFromConfig(cfg)

	out, err := stsClient.GetCallerIdentity(t.Context(), &sts.GetCallerIdentityInput{})
	if err != nil {
		t.Fatalf("failed to get caller identity: %v", err)
	}

	if out.Account == nil {
		t.Fatal("aws sts get-caller-identity returned nil account ID")
	}

	accountID := *out.Account
	region := cfg.Region

	return &AccTestEnv{
		Region:     region,
		AccountID:  accountID,
		BucketName: fmt.Sprintf("%s-%s-%s", BucketPrefix, accountID, region),
	}
}

// DefaultVPCInfo describes the default VPC in the current AWS account and region.
//
// It includes the VPC ID and a sorted list of subnet IDs and is primarily used
// by acceptance tests that require networking prerequisites (for example,
// AppStream resources that must be deployed into a VPC).
type DefaultVPCInfo struct {
	// VpcID is the identifier of the default VPC.
	VpcID string
	// SubnetIDs contains the IDs of all subnets associated with the default VPC.
	// The slice is sorted to ensure deterministic test behavior.
	SubnetIDs []string
}

// DefaultVPCInfo retrieves information about the default VPC in the AWS account
// and region associated with the AccTestEnv.
//
// The function:
//   - Uses the region stored in the AccTestEnv
//   - Requires that a default VPC exists
//   - Requires that the default VPC has at least two subnets
//
// The test fails immediately if these conditions are not met, as they are
// mandatory for many AppStream acceptance tests.
func (ctx *AccTestEnv) DefaultVPCInfo(t *testing.T) *DefaultVPCInfo {
	t.Helper()

	cfg, err := config.LoadDefaultConfig(t.Context(), config.WithRegion(ctx.Region))
	if err != nil {
		t.Fatalf("failed to load aws config: %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	// find default VPC
	vpcsOut, err := ec2Client.DescribeVpcs(t.Context(), &ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("isDefault"),
				Values: []string{"true"},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to describe vpcs: %v", err)
	}

	if len(vpcsOut.Vpcs) == 0 {
		t.Fatal("no default vpc found")
	}

	vpcID := aws.ToString(vpcsOut.Vpcs[0].VpcId)

	// find subnets in that VPC
	subnetsOut, err := ec2Client.DescribeSubnets(t.Context(), &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to describe subnets: %v", err)
	}

	var subnetIDs []string
	for _, s := range subnetsOut.Subnets {
		if s.SubnetId != nil {
			subnetIDs = append(subnetIDs, *s.SubnetId)
		}
	}

	sort.Strings(subnetIDs)

	if len(subnetIDs) < 2 {
		t.Fatalf("default vpc %s has fewer than 2 subnets (found %d)", vpcID, len(subnetIDs))
	}

	return &DefaultVPCInfo{VpcID: vpcID, SubnetIDs: subnetIDs}
}
