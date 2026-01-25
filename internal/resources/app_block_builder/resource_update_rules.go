// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import (
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

type appBlockBuilderUpdateMode int

const (
	appBlockBuilderUpdateAllowedRunning appBlockBuilderUpdateMode = iota
	appBlockBuilderUpdateRequiresStop
)

func updateMode(state, plan resourceModel) appBlockBuilderUpdateMode {
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateAppBlockBuilder.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_UpdateAppBlockBuilder.html

	if util.Changed(state.InstanceType, plan.InstanceType) {
		return appBlockBuilderUpdateRequiresStop
	}

	if util.Changed(state.Platform, plan.Platform) {
		return appBlockBuilderUpdateRequiresStop
	}

	if util.Changed(state.VPCConfig, plan.VPCConfig) {
		return appBlockBuilderUpdateRequiresStop
	}

	if util.Changed(state.IAMRoleARN, plan.IAMRoleARN) {
		return appBlockBuilderUpdateRequiresStop
	}

	if util.Changed(state.EnableDefaultInternetAccess, plan.EnableDefaultInternetAccess) {
		return appBlockBuilderUpdateRequiresStop
	}

	if util.Changed(state.AccessEndpoints, plan.AccessEndpoints) {
		return appBlockBuilderUpdateRequiresStop
	}

	return appBlockBuilderUpdateAllowedRunning
}
