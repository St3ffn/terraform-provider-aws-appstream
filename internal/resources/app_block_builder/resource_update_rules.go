// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

type appBlockBuilderUpdateMode int

const (
	appBlockBuilderUpdateAllowedRunning appBlockBuilderUpdateMode = iota
	appBlockBuilderUpdateRequiresStop
)

func updateMode(diff resourceDiff) appBlockBuilderUpdateMode {
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateAppBlockBuilder.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_UpdateAppBlockBuilder.html

	if diff.InstanceType.IsChanged() {
		return appBlockBuilderUpdateRequiresStop
	}

	if diff.Platform.IsChanged() {
		return appBlockBuilderUpdateRequiresStop
	}

	if diff.VPCConfig.IsChanged() {
		return appBlockBuilderUpdateRequiresStop
	}

	if diff.IAMRoleARN.IsChanged() {
		return appBlockBuilderUpdateRequiresStop
	}

	if diff.DisableIMDSV1.IsChanged() {
		return appBlockBuilderUpdateRequiresStop
	}

	if diff.EnableDefaultInternetAccess.IsChanged() {
		return appBlockBuilderUpdateRequiresStop
	}

	if diff.AccessEndpoints.IsChanged() {
		return appBlockBuilderUpdateRequiresStop
	}

	return appBlockBuilderUpdateAllowedRunning
}
