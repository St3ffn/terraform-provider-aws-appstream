// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"context"
	"fmt"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan computes tags_all and performs a plan-time restart risk check.
// If the diff requires a stop/start update, update_behavior is AUTO_STOP_START,
// and current state is RUNNING, it emits a warning before apply.
func (r *resource) ModifyPlan(ctx context.Context, req tfresource.ModifyPlanRequest, resp *tfresource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state resourceModel
	hasState := !req.State.Raw.IsNull()
	if hasState {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan.TagsAll = r.tags.EffectiveTagsForPlan(plan.Tags)

	if hasState {
		diff := newResourceDiff(state, plan)
		mode := updateMode(plan, diff)
		behavior := updateBehaviorFromPlan(plan.UpdateBehavior)
		desired := desiredStateFromPlan(plan.DesiredState)

		currentStateKnown := !state.State.IsNull() && !state.State.IsUnknown()
		currentStateRunning := currentStateKnown && wasRunning(
			awstypes.FleetState(state.State.ValueString()),
		)

		if mode == fleetUpdateRequiresStop &&
			behavior == updateBehaviorAutoStopStart &&
			desired != desiredStateStopped &&
			currentStateRunning {
			resp.Diagnostics.AddWarning(
				"AWS AppStream Fleet May Be Restarted",
				fmt.Sprintf(
					"Updating fleet %q requires the fleet to be stopped and update_behavior is %q. "+
						"Because the current state is %q, the provider plans to stop and start the fleet during apply.",
					plan.Name.ValueString(),
					updateBehaviorAutoStopStart.String(),
					state.State.ValueString(),
				),
			)
		}
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}
