// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
)

type fleetUpdateMode int

const (
	fleetUpdateAllowedRunning fleetUpdateMode = iota
	fleetUpdateRequiresStop
	fleetUpdateForbidden
)

func updateMode(plan resourceModel, diff resourceDiff) fleetUpdateMode {
	switch awstypes.FleetType(plan.FleetType.ValueString()) {
	case awstypes.FleetTypeAlwaysOn, awstypes.FleetTypeOnDemand:
		return updateModeAlwaysOnOnDemandFleet(diff)
	case awstypes.FleetTypeElastic:
		return updateModeElasticFleet(diff)
	default:
		return fleetUpdateForbidden
	}
}

func updateModeAlwaysOnOnDemandFleet(diff resourceDiff) fleetUpdateMode {
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateFleet.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_UpdateFleet.html

	if diff.InstanceType.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.VPCConfig.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.MaxUserDurationInSeconds.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.DisableIMDSV1.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.EnableDefaultInternetAccess.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.DomainJoinInfo.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.IAMRoleARN.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.StreamView.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.Platform.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.MaxSessionsPerInstance.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.USBDeviceFilterStrings.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.RootVolumeConfig.IsChanged() {
		return fleetUpdateRequiresStop
	}

	return fleetUpdateAllowedRunning
}

func updateModeElasticFleet(diff resourceDiff) fleetUpdateMode {
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateFleet.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_UpdateFleet.html

	if diff.ImageName.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.ImageARN.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.InstanceType.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.VPCConfig.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.MaxUserDurationInSeconds.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.DisableIMDSV1.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.EnableDefaultInternetAccess.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.IAMRoleARN.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.StreamView.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.Platform.IsChanged() {
		return fleetUpdateRequiresStop
	}

	if diff.RootVolumeConfig.IsChanged() {
		return fleetUpdateRequiresStop
	}

	return fleetUpdateAllowedRunning
}
