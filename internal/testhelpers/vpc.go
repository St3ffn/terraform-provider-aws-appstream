// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// DefaultVPCInfo describes the default VPC in the current AWS account and region,
// including its VPC ID and associated subnet IDs.
//
// This structure is primarily used by acceptance tests that require networking
// prerequisites (for example, AppStream resources that must be placed into a VPC).
type DefaultVPCInfo struct {
	// VpcID is the identifier of the default VPC.
	VpcID string
	// SubnetIDs contains the IDs of all subnets associated with the default VPC.
	// The list is sorted lexicographically to ensure deterministic test behavior.
	SubnetIDs []string
}

// GetDefaultVPCInfo retrieves information about the default VPC in the current AWS
// account and region, including its VPC ID and subnet IDs.
//
// This helper is intended for use in Terraform acceptance tests only. If acceptance
// tests are not enabled (TF_ACC is not set), the test is skipped automatically.
//
// The function requires that a default VPC exists and that it contains at least
// two subnets, which is a common requirement for AppStream and other AWS services.
// An error is returned if these conditions are not met.
//
// AWS credentials and region configuration are resolved using the default AWS SDK
// configuration chain.
func GetDefaultVPCInfo(t *testing.T) (*DefaultVPCInfo, error) {
	if !IsAccTest() {
		t.Skip("Skipping acceptance test unless TF_ACC is set")
	}

	t.Helper()

	cfg, err := config.LoadDefaultConfig(t.Context())
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
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
		return nil, fmt.Errorf("describing VPCs: %w", err)
	}

	if len(vpcsOut.Vpcs) == 0 {
		return nil, errors.New("no default VPC found; acceptance tests require a default VPC or at least two subnets")
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
		return nil, fmt.Errorf("describing subnets: %w", err)
	}

	var subnetIDs []string
	for _, s := range subnetsOut.Subnets {
		if s.SubnetId != nil {
			subnetIDs = append(subnetIDs, *s.SubnetId)
		}
	}

	sort.Strings(subnetIDs)

	if len(subnetIDs) < 2 {
		return nil, fmt.Errorf(
			"default VPC %s has fewer than 2 subnets (found %d)",
			vpcID, len(subnetIDs),
		)
	}

	return &DefaultVPCInfo{
		VpcID:     vpcID,
		SubnetIDs: subnetIDs,
	}, nil
}
