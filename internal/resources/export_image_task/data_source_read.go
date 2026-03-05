// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package export_image_task

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Read fetches the current remote object and refreshes Terraform state from AWS.
// When the object no longer exists remotely, the resource is removed from state
// to converge Terraform with external deletions.
func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if config.TaskID.IsNull() || config.TaskID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read export image task because task_id must be set and known.",
		)
		return
	}

	taskID := config.TaskID.ValueString()

	out, err := ds.appstreamClient.GetExportImageTask(ctx, &awsappstream.GetExportImageTaskInput{
		TaskId: aws.String(taskID),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Export Image Task Not Found",
				fmt.Sprintf("No export image task %q was found.", taskID),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Export Image Task",
			fmt.Sprintf("Could not read export image task %q: %v", taskID, err),
		)
		return
	}

	if out.ExportImageTask == nil {
		resp.Diagnostics.AddError(
			"AWS AppStream Export Image Task Not Found",
			fmt.Sprintf("No export image task %q was found.", taskID),
		)
		return
	}

	task := out.ExportImageTask

	state := &dataSourceModel{
		ID:                types.StringValue(aws.ToString(task.TaskId)),
		TaskID:            types.StringValue(aws.ToString(task.TaskId)),
		ImageARN:          util.StringOrNull(task.ImageArn),
		AmiName:           util.StringOrNull(task.AmiName),
		AmiDescription:    util.StringOrNull(task.AmiDescription),
		AmiID:             util.StringOrNull(task.AmiId),
		State:             types.StringValue(string(task.State)),
		CreatedDate:       util.StringFromTime(task.CreatedDate),
		TagSpecifications: util.MapStringOrNull(ctx, task.TagSpecifications, &resp.Diagnostics),
		ErrorDetails:      flattenErrorDetailsData(ctx, task.ErrorDetails, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
