// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package export_image_tasks

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	filters := expandFilters(ctx, config.Filters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var tasks []awstypes.ExportImageTask
	var nextToken *string

	for {
		out, err := ds.appstreamClient.ListExportImageTasks(ctx, &awsappstream.ListExportImageTasksInput{
			Filters:   filters,
			NextToken: nextToken,
		})
		if err != nil {
			if util.IsContextCanceled(err) {
				return
			}

			resp.Diagnostics.AddError(
				"Error Reading AWS AppStream Export Image Tasks",
				fmt.Sprintf("Could not read export image tasks: %v", err),
			)
			return
		}

		tasks = append(tasks, out.ExportImageTasks...)

		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	state := &dataSourceModel{
		Filters:          config.Filters,
		ExportImageTasks: flattenExportImageTasksData(ctx, tasks, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
