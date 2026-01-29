// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addDiagnostics(state, &resp.Diagnostics, diagnosticRead)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readAssociationImageBuilderSoftware(ctx, state)
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

func (r *resource) readAssociationImageBuilderSoftware(ctx context.Context, prior model) (*model, diag.Diagnostics) {
	var diags diag.Diagnostics

	imageBuilderARN := prior.ImageBuilderARN.ValueString()

	out, err := r.appstreamClient.DescribeSoftwareAssociations(ctx, &awsappstream.DescribeSoftwareAssociationsInput{
		AssociatedResource: aws.String(imageBuilderARN),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream Image Builder Software Association",
			fmt.Sprintf("Could not read association of software for image builder %q: %v", imageBuilderARN, err),
		)
		return nil, diags
	}

	if len(out.SoftwareAssociations) == 0 {
		return nil, diags
	}

	filtered := filterSoftwareAssociations(ctx, prior.SoftwareNames, out.SoftwareAssociations, &diags)
	if diags.HasError() {
		return nil, diags
	}

	state := &model{
		ID:              prior.ImageBuilderARN,
		ImageBuilderARN: prior.ImageBuilderARN,
		SoftwareNames:   prior.SoftwareNames,
		Deploy:          prior.Deploy,
		Associations:    flattenAssociations(ctx, filtered, &diags),
	}

	return state, diags
}

func filterSoftwareAssociations(
	ctx context.Context, prior types.Set, awsAssociations []awstypes.SoftwareAssociations, diags *diag.Diagnostics,
) []awstypes.SoftwareAssociations {

	// if terraform does not know what it manages, do not adopt anything.
	if prior.IsNull() || prior.IsUnknown() {
		return nil
	}

	// decode managed software names
	var names []string
	diags.Append(prior.ElementsAs(ctx, &names, false)...)
	if diags.HasError() {
		return nil
	}

	index := make(map[string]struct{}, len(names))
	for _, n := range names {
		index[n] = struct{}{}
	}

	out := make([]awstypes.SoftwareAssociations, 0, len(index))
	for _, assoc := range awsAssociations {
		if assoc.SoftwareName == nil {
			continue
		}

		if _, ok := index[*assoc.SoftwareName]; ok {
			out = append(out, assoc)
		}
	}

	return out
}
