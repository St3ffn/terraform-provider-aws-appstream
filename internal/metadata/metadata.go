// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

// Metadata contains provider-scoped AWS clients and defaults shared with
// resources and data sources during Configure.
type Metadata struct {
	Appstream   *awsappstream.Client
	Tagging     *awstaggingapi.Client
	DefaultTags map[string]string
}

// NewMetadata builds provider metadata from AWS config and default tags.
func NewMetadata(awscfg aws.Config, defaultTags map[string]string) *Metadata {
	return &Metadata{
		Appstream:   awsappstream.NewFromConfig(awscfg),
		Tagging:     awstaggingapi.NewFromConfig(awscfg),
		DefaultTags: defaultTags,
	}
}
