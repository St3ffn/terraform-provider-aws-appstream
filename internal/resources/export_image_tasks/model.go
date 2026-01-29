// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package export_image_tasks

import "github.com/hashicorp/terraform-plugin-framework/types"

type model struct {
	// Filters are optional criteria used to narrow down the list of export image tasks (optional).
	Filters types.Set `tfsdk:"filters"`
	// ExportImageTasks is the list of export image tasks that match the specified filters (computed).
	ExportImageTasks types.Set `tfsdk:"export_image_tasks"`
}

type filterModel struct {
	// Name is the name of the filter applied when listing export image tasks (required).
	Name types.String `tfsdk:"name"`
	// Values are the values for the specified filter name (required).
	Values types.Set `tfsdk:"values"`
}

type exportImageTaskModel struct {
	// TaskID is the unique identifier of the export image task (computed).
	TaskID types.String `tfsdk:"task_id"`
	// ImageArn is the ARN of the AppStream image being exported (computed).
	ImageArn types.String `tfsdk:"image_arn"`
	// AmiName is the name of the EC2 AMI created by the export task (computed).
	AmiName types.String `tfsdk:"ami_name"`
	// AmiDescription is the description applied to the exported EC2 AMI (computed).
	AmiDescription types.String `tfsdk:"ami_description"`
	// AmiID is the ID of the EC2 AMI created by the export task (computed).
	AmiID types.String `tfsdk:"ami_id"`
	// State is the current state of the export image task (computed).
	State types.String `tfsdk:"state"`
	// CreatedDate is the timestamp when the export image task was created (computed).
	CreatedDate types.String `tfsdk:"created_date"`
	// TagSpecifications are the tags applied to the exported EC2 AMI (computed).
	TagSpecifications types.Map `tfsdk:"tag_specifications"`
	// ErrorDetails is the list of errors reported by AWS for the export image task (computed).
	ErrorDetails types.Set `tfsdk:"error_details"`
}

type errorDetailModel struct {
	// ErrorCode is the error code reported by AWS for the export task (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message reported by AWS (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
}
