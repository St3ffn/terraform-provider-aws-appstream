// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func storageConnectorAttributesToDelete(
	ctx context.Context, priorSet types.Set, planSet types.Set, diags *diag.Diagnostics,
) []awstypes.StackAttribute {

	removedConnectorTypes := make(map[awstypes.StorageConnectorType]struct{})
	attrsToDelete := make([]awstypes.StackAttribute, 0)

	if !priorSet.IsNull() && !priorSet.IsUnknown() {
		var prior []storageConnectorModel
		diags.Append(priorSet.ElementsAs(ctx, &prior, false)...)
		if diags.HasError() {
			return nil
		}

		for _, c := range prior {
			t := awstypes.StorageConnectorType(c.ConnectorType.ValueString())
			removedConnectorTypes[t] = struct{}{}
		}
	}

	if !planSet.IsNull() && !planSet.IsUnknown() {
		var desired []storageConnectorModel
		diags.Append(planSet.ElementsAs(ctx, &desired, false)...)
		if diags.HasError() {
			return nil
		}

		for _, c := range desired {
			t := awstypes.StorageConnectorType(c.ConnectorType.ValueString())
			delete(removedConnectorTypes, t)
		}
	}

	for t := range removedConnectorTypes {
		if attr, ok := storageConnectorDeleteAttribute(t); ok {
			attrsToDelete = append(attrsToDelete, attr)
		}
	}

	return attrsToDelete
}

func storageConnectorDeleteAttribute(t awstypes.StorageConnectorType) (awstypes.StackAttribute, bool) {
	switch t {
	case awstypes.StorageConnectorTypeHomefolders:
		return awstypes.StackAttributeStorageConnectorHomefolders, true
	case awstypes.StorageConnectorTypeGoogleDrive:
		return awstypes.StackAttributeStorageConnectorGoogleDrive, true
	case awstypes.StorageConnectorTypeOneDrive:
		return awstypes.StackAttributeStorageConnectorOneDrive, true
	default:
		return "", false
	}
}
