// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

type fleetUpdateMode int

const (
	fleetUpdateAllowedRunning fleetUpdateMode = iota
	fleetUpdateRequiresStop
	fleetUpdateForbidden
)

func updateMode(state, plan resourceModel) fleetUpdateMode {
	if plan.State.IsUnknown() {
		return fleetUpdateForbidden
	}

	switch awstypes.FleetType(plan.FleetType.ValueString()) {
	case awstypes.FleetTypeAlwaysOn, awstypes.FleetTypeOnDemand:
		return updateModeAlwaysOnOnDemandFleet(state, plan)
	case awstypes.FleetTypeElastic:
		return updateModeElasticFleet(state, plan)
	default:
		return fleetUpdateForbidden
	}
}

func updateModeAlwaysOnOnDemandFleet(state, plan resourceModel) fleetUpdateMode {
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateFleet.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_UpdateFleet.html

	if util.Changed(state.InstanceType, plan.InstanceType) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.VPCConfig, plan.VPCConfig) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.MaxUserDurationInSeconds, plan.MaxUserDurationInSeconds) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.EnableDefaultInternetAccess, plan.EnableDefaultInternetAccess) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.DomainJoinInfo, plan.DomainJoinInfo) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.IAMRoleARN, plan.IAMRoleARN) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.StreamView, plan.StreamView) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.Platform, plan.Platform) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.MaxSessionsPerInstance, plan.MaxSessionsPerInstance) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.USBDeviceFilterStrings, plan.USBDeviceFilterStrings) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.RootVolumeConfig, plan.RootVolumeConfig) {
		return fleetUpdateRequiresStop
	}

	return fleetUpdateAllowedRunning
}

func updateModeElasticFleet(state, plan resourceModel) fleetUpdateMode {
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateFleet.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_UpdateFleet.html

	if util.Changed(state.ImageName, plan.ImageName) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.ImageARN, plan.ImageARN) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.InstanceType, plan.InstanceType) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.VPCConfig, plan.VPCConfig) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.MaxUserDurationInSeconds, plan.MaxUserDurationInSeconds) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.EnableDefaultInternetAccess, plan.EnableDefaultInternetAccess) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.IAMRoleARN, plan.IAMRoleARN) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.StreamView, plan.StreamView) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.Platform, plan.Platform) {
		return fleetUpdateRequiresStop
	}

	if util.Changed(state.RootVolumeConfig, plan.RootVolumeConfig) {
		return fleetUpdateRequiresStop
	}

	return fleetUpdateAllowedRunning
}
