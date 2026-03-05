// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import "github.com/hashicorp/terraform-plugin-framework/types"

type desiredState string

// String returns the string representation of the current value.
func (s desiredState) String() string {
	return string(s)
}

const (
	desiredStateInherit desiredState = "INHERIT"
	desiredStateRunning desiredState = "RUNNING"
	desiredStateStopped desiredState = "STOPPED"
)

func desiredStateFromPlan(v types.String) desiredState {
	if v.IsNull() || v.IsUnknown() {
		return desiredStateInherit
	}

	return desiredState(v.ValueString())
}

type updateBehavior string

// String returns the string representation of the current value.
func (b updateBehavior) String() string {
	return string(b)
}

const (
	updateBehaviorAutoStopStart updateBehavior = "AUTO_STOP_START"
	updateBehaviorFailIfRunning updateBehavior = "FAIL_IF_RUNNING"
)

func updateBehaviorFromPlan(v types.String) updateBehavior {
	if v.IsNull() || v.IsUnknown() {
		return updateBehaviorAutoStopStart
	}

	return updateBehavior(v.ValueString())
}
