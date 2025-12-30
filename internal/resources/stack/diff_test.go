// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func testStorageConnector(connectorType string) storageConnectorModel {
	return storageConnectorModel{
		ConnectorType:              types.StringValue(connectorType),
		ResourceIdentifier:         types.StringNull(),
		Domains:                    types.SetNull(types.StringType),
		DomainsRequireAdminConsent: types.SetNull(types.StringType),
	}
}

func storageConnectorSet(ctx context.Context, t *testing.T, models []storageConnectorModel) types.Set {
	t.Helper()

	if len(models) == 0 {
		return types.SetNull(storageConnectorObjectType)
	}

	setVal, diags := types.SetValueFrom(ctx, storageConnectorObjectType, models)
	require.False(t, diags.HasError(), "failed_to_build_storage_connector_set: %v", diags)

	return setVal
}

func TestStorageConnectorAttributesToDelete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prior     []storageConnectorModel
		plan      []storageConnectorModel
		wantAttrs []awstypes.StackAttribute
	}{
		{
			name:      "no_prior_no_plan",
			prior:     nil,
			plan:      nil,
			wantAttrs: nil,
		},
		{
			name: "prior_only_delete_all",
			prior: []storageConnectorModel{
				testStorageConnector("HOMEFOLDERS"),
				testStorageConnector("GOOGLE_DRIVE"),
			},
			plan: nil,
			wantAttrs: []awstypes.StackAttribute{
				awstypes.StackAttributeStorageConnectorHomefolders,
				awstypes.StackAttributeStorageConnectorGoogleDrive,
			},
		},
		{
			name: "prior_and_plan_identical_no_deletes",
			prior: []storageConnectorModel{
				testStorageConnector("HOMEFOLDERS"),
			},
			plan: []storageConnectorModel{
				testStorageConnector("HOMEFOLDERS"),
			},
			wantAttrs: nil,
		},
		{
			name: "remove_one_connector",
			prior: []storageConnectorModel{
				testStorageConnector("HOMEFOLDERS"),
				testStorageConnector("GOOGLE_DRIVE"),
			},
			plan: []storageConnectorModel{
				testStorageConnector("HOMEFOLDERS"),
			},
			wantAttrs: []awstypes.StackAttribute{
				awstypes.StackAttributeStorageConnectorGoogleDrive,
			},
		},
		{
			name: "unknown_connector_type_ignored",
			prior: []storageConnectorModel{
				testStorageConnector("HOMEFOLDERS"),
				testStorageConnector("UNKNOWN"),
			},
			plan: []storageConnectorModel{
				testStorageConnector("HOMEFOLDERS"),
			},
			wantAttrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			priorSet := storageConnectorSet(ctx, t, tt.prior)
			planSet := storageConnectorSet(ctx, t, tt.plan)

			attrs := storageConnectorAttributesToDelete(ctx, priorSet, planSet, &diags)

			require.False(t, diags.HasError(), "unexpected_diagnostics_for_test_%s: %v", tt.name, diags)

			require.ElementsMatch(t, tt.wantAttrs, attrs, "unexpected_attributes_to_delete_for_test_%s", tt.name)
		})
	}
}

func TestStorageConnectorDeleteAttribute(t *testing.T) {
	tests := []struct {
		name      string
		input     awstypes.StorageConnectorType
		wantAttr  awstypes.StackAttribute
		wantFound bool
	}{
		{
			name:      "HOMEFOLDERS",
			input:     awstypes.StorageConnectorTypeHomefolders,
			wantAttr:  awstypes.StackAttributeStorageConnectorHomefolders,
			wantFound: true,
		},
		{
			name:      "GOOGLE_DRIVE",
			input:     awstypes.StorageConnectorTypeGoogleDrive,
			wantAttr:  awstypes.StackAttributeStorageConnectorGoogleDrive,
			wantFound: true,
		},
		{
			name:      "ONE_DRIVE",
			input:     awstypes.StorageConnectorTypeOneDrive,
			wantAttr:  awstypes.StackAttributeStorageConnectorOneDrive,
			wantFound: true,
		},
		{
			name:      "unknown_type",
			input:     awstypes.StorageConnectorType("SOMETHING_ELSE"),
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr, found := storageConnectorDeleteAttribute(tt.input)

			require.Equal(t, tt.wantFound, found, "unexpected found flag for connector type %q", tt.input)

			if tt.wantFound {
				require.Equal(
					t, tt.wantAttr, attr, "unexpected stack attribute for connector type %q", tt.input,
				)
			}
		})
	}
}
