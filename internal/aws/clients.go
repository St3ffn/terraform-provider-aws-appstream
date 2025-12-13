// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type Clients struct {
	AppStream *awsappstream.Client
	Tagging   *awstaggingapi.Client
}

func NewClients(awscfg aws.Config) *Clients {
	return &Clients{
		AppStream: awsappstream.NewFromConfig(awscfg),
		Tagging:   awstaggingapi.NewFromConfig(awscfg),
	}
}
