// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *associateApplicationEntitlementResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// no-op: all attributes require replacement
}
