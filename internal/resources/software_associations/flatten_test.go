// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package software_associations

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenSoftwareAssociations(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.SoftwareAssociations
		want types.Set
	}{
		{
			name: "empty_input",
			in:   nil,
			want: types.SetNull(softwareAssociationObjectType),
		},
		{
			name: "single_association_no_errors",
			in: []awstypes.SoftwareAssociations{
				{
					SoftwareName: aws.String("Microsoft_Office_2021_LTSC_Professional_Plus_64Bit"),
					Status:       awstypes.SoftwareDeploymentStatusInstalled,
				},
			},
			want: types.SetValueMust(
				softwareAssociationObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						softwareAssociationObjectType.AttrTypes,
						map[string]attr.Value{
							"software_name": types.StringValue("Microsoft_Office_2021_LTSC_Professional_Plus_64Bit"),
							"status":        types.StringValue("INSTALLED"),
							"deployment_errors": types.SetNull(
								deploymentErrorObjectType,
							),
						},
					),
				},
			),
		},
		{
			name: "single_association_with_errors",
			in: []awstypes.SoftwareAssociations{
				{
					SoftwareName: aws.String("Microsoft_Project_2024_Professional_64Bit"),
					Status:       awstypes.SoftwareDeploymentStatusFailedToInstall,
					DeploymentError: []awstypes.ErrorDetails{
						{
							ErrorCode:    aws.String("INSTALL_FAILED"),
							ErrorMessage: aws.String("Installation failed"),
						},
					},
				},
			},
			want: types.SetValueMust(
				softwareAssociationObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						softwareAssociationObjectType.AttrTypes,
						map[string]attr.Value{
							"software_name": types.StringValue("Microsoft_Project_2024_Professional_64Bit"),
							"status":        types.StringValue("FAILED_TO_INSTALL"),
							"deployment_errors": types.SetValueMust(
								deploymentErrorObjectType,
								[]attr.Value{
									types.ObjectValueMust(
										deploymentErrorObjectType.AttrTypes,
										map[string]attr.Value{
											"error_code":    types.StringValue("INSTALL_FAILED"),
											"error_message": types.StringValue("Installation failed"),
										},
									),
								},
							),
						},
					),
				},
			),
		},
		{
			name: "multiple_associations",
			in: []awstypes.SoftwareAssociations{
				{
					SoftwareName: aws.String("Microsoft_Visio_2024_LTSC_Standard_64Bit"),
					Status:       awstypes.SoftwareDeploymentStatusInstalled,
				},
				{
					SoftwareName: aws.String("Microsoft_Project_2021_Standard_32Bit"),
					Status:       awstypes.SoftwareDeploymentStatusPendingInstallation,
				},
			},
			want: types.SetValueMust(
				softwareAssociationObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						softwareAssociationObjectType.AttrTypes,
						map[string]attr.Value{
							"software_name": types.StringValue("Microsoft_Visio_2024_LTSC_Standard_64Bit"),
							"status":        types.StringValue("INSTALLED"),
							"deployment_errors": types.SetNull(
								deploymentErrorObjectType,
							),
						},
					),
					types.ObjectValueMust(
						softwareAssociationObjectType.AttrTypes,
						map[string]attr.Value{
							"software_name": types.StringValue("Microsoft_Project_2021_Standard_32Bit"),
							"status":        types.StringValue("PENDING_INSTALLATION"),
							"deployment_errors": types.SetNull(
								deploymentErrorObjectType,
							),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenSoftwareAssociations(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenSoftwareAssociationDeploymentErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.ErrorDetails
		want types.Set
	}{
		{
			name: "empty_input",
			in:   nil,
			want: types.SetNull(deploymentErrorObjectType),
		},
		{
			name: "single_error",
			in: []awstypes.ErrorDetails{
				{
					ErrorCode:    aws.String("ERROR_CODE"),
					ErrorMessage: aws.String("something went wrong"),
				},
			},
			want: types.SetValueMust(
				deploymentErrorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						deploymentErrorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":    types.StringValue("ERROR_CODE"),
							"error_message": types.StringValue("something went wrong"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenSoftwareAssociationDeploymentErrors(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
