// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/require"
)

func TestFlattenAssociations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name           string
		in             []awstypes.SoftwareAssociations
		expectNull     bool
		expectElements int
	}{
		{
			name:       "empty_input_returns_null",
			in:         nil,
			expectNull: true,
		},
		{
			name: "single_association_no_errors",
			in: []awstypes.SoftwareAssociations{
				{
					SoftwareName: aws.String("Microsoft_Office_2024_LTSC_Professional_Plus_64Bit"),
					Status:       awstypes.SoftwareDeploymentStatusPendingInstallation,
				},
			},
			expectElements: 1,
		},
		{
			name: "single_association_with_error",
			in: []awstypes.SoftwareAssociations{
				{
					SoftwareName: aws.String("Microsoft_Project_2024_Professional_64Bit"),
					Status:       awstypes.SoftwareDeploymentStatusFailedToInstall,
					DeploymentError: []awstypes.ErrorDetails{
						{
							ErrorCode:    aws.String("INTERNAL_ERROR"),
							ErrorMessage: aws.String("boom"),
						},
					},
				},
			},
			expectElements: 1,
		},
		{
			name: "multiple_associations",
			in: []awstypes.SoftwareAssociations{
				{
					SoftwareName: aws.String("SoftwareA"),
					Status:       awstypes.SoftwareDeploymentStatusInstalled,
				},
				{
					SoftwareName: aws.String("SoftwareB"),
					Status:       awstypes.SoftwareDeploymentStatusPendingInstallation,
				},
			},
			expectElements: 2,
		},
		{
			name: "nil_software_name_is_ignored",
			in: []awstypes.SoftwareAssociations{
				{
					SoftwareName: nil,
					Status:       awstypes.SoftwareDeploymentStatusInstalled,
				},
			},
			expectElements: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			out := flattenAssociations(ctx, tt.in, &diags)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			if tt.expectNull {
				require.True(t, out.IsNull(), "expected null map")
				return
			}

			require.False(t, out.IsNull(), "expected non-null map")

			var m map[string]associationModel
			diags = out.ElementsAs(ctx, &m, false)
			require.False(t, diags.HasError(), "failed to decode map elements")

			require.Len(t, m, tt.expectElements, "unexpected number of association entries")

			// spot-check contents for error case
			if tt.name == "single_association_with_error" {
				assoc, ok := m["Microsoft_Project_2024_Professional_64Bit"]
				require.True(t, ok, "expected association key not found")

				require.Equal(
					t,
					string(awstypes.SoftwareDeploymentStatusFailedToInstall),
					assoc.Status.ValueString(),
					"unexpected status value",
				)

				var errs []deploymentErrorModel
				diags = assoc.DeploymentErrors.ElementsAs(ctx, &errs, false)
				require.False(t, diags.HasError(), "failed to decode deployment errors")

				require.Len(t, errs, 1, "expected exactly one deployment error")
				require.Equal(t, "INTERNAL_ERROR", errs[0].ErrorCode.ValueString())
				require.Equal(t, "boom", errs[0].ErrorMessage.ValueString())
			}
		})
	}
}
