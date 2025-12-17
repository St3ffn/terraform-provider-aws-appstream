// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type awsClients struct {
	AppStream *awsappstream.Client
	Tagging   *awstaggingapi.Client
}

func newAWSClients(awscfg aws.Config) *awsClients {
	return &awsClients{
		AppStream: awsappstream.NewFromConfig(awscfg),
		Tagging:   awstaggingapi.NewFromConfig(awscfg),
	}
}
