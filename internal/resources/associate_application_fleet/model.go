// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_fleet

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type model struct {
	// ID is a synthetic identifier composed of "<fleet_name>|<application_arn>".
	ID types.String `tfsdk:"id"`
	// FleetName is the name of the AppStream fleet to be associated with the application (required).
	FleetName types.String `tfsdk:"fleet_name"`
	// ApplicationARN is the ARN of the AppStream application to associate with the fleet (required).
	ApplicationARN types.String `tfsdk:"application_arn"`
}
