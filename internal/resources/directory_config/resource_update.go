// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/path"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Update(ctx context.Context, req tfresource.UpdateRequest, resp *tfresource.UpdateResponse) {
	var plan resourceModel
	var state resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if plan.DirectoryName.IsNull() || plan.DirectoryName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update directory config because directory_name must be known.",
		)
		return
	}

	diff := newResourceDiff(state, plan)
	name := plan.DirectoryName.ValueString()

	if diff.CertificateBasedAuthProperties.IsCleared() &&
		!state.CertificateBasedAuthProperties.IsNull() &&
		!state.CertificateBasedAuthProperties.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("certificate_based_auth_properties"),
			"Cannot unset certificate_based_auth_properties",
			"AWS UpdateDirectoryConfig does not support removing certificate_based_auth_properties once set. "+
				"Use certificate_based_auth_properties.status = \"DISABLED\" instead.",
		)
		return
	}

	if diff.HasRemoteChanges() {
		input := &awsappstream.UpdateDirectoryConfigInput{
			DirectoryName: aws.String(name),
		}

		if diff.CertificateBasedAuthProperties.IsUpdated() {
			input.CertificateBasedAuthProperties = expandCertificateBasedAuthProperties(
				ctx, plan.CertificateBasedAuthProperties, &resp.Diagnostics,
			)
		}

		if diff.OrganizationalUnitDistinguishedNames.IsChanged() {
			input.OrganizationalUnitDistinguishedNames = util.ExpandStringSetOrNil(
				ctx, plan.OrganizationalUnitDistinguishedNames, &resp.Diagnostics,
			)
		}

		if diff.ServiceAccountCredentials.IsChanged() && !plan.ServiceAccountCredentials.IsNull() {
			input.ServiceAccountCredentials = expandServiceAccountCredentials(
				ctx, plan.ServiceAccountCredentials, &resp.Diagnostics,
			)
		}

		if resp.Diagnostics.HasError() {
			return
		}

		_, err := r.appstreamClient.UpdateDirectoryConfig(ctx, input)
		if err != nil {
			if util.IsContextCanceled(err) {
				return
			}

			if util.IsAppStreamNotFound(err) {
				resp.State.RemoveResource(ctx)
				return
			}

			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Directory Config",
				fmt.Sprintf("Could not update directory config %q: %v", name, err),
			)
			return
		}
	}

	newState, diags := r.readDirectoryConfig(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if ctx.Err() != nil {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
