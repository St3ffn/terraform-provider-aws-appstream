// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type Metadata struct {
	Appstream   *awsappstream.Client
	Tagging     *awstaggingapi.Client
	DefaultTags map[string]string
}

func NewMetadata(awscfg aws.Config, defaultTags map[string]string) *Metadata {
	return &Metadata{
		Appstream:   awsappstream.NewFromConfig(awscfg),
		Tagging:     awstaggingapi.NewFromConfig(awscfg),
		DefaultTags: defaultTags,
	}
}
