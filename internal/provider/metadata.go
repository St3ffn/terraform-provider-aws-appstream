// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type metadata struct {
	appstream   *awsappstream.Client
	tagging     *awstaggingapi.Client
	defaultTags map[string]string
}

func newMetadata(awscfg aws.Config, defaultTags map[string]string) *metadata {
	return &metadata{
		appstream:   awsappstream.NewFromConfig(awscfg),
		tagging:     awstaggingapi.NewFromConfig(awscfg),
		defaultTags: defaultTags,
	}
}

//lint:ignore U1000 applyDefaultTags will be used in a future change
func (m *metadata) applyDefaultTags(tags map[string]string) map[string]string {
	out := make(map[string]string)

	for k, v := range m.defaultTags {
		out[k] = v
	}
	for k, v := range tags {
		out[k] = v
	}
	return out
}
