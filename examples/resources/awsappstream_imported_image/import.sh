# Import uses image name only.
# Note: create-time-only attributes (iam_role_arn, source_ami_id, agent_software_version,
# runtime_validation_config, app_catalog_config) are not returned by DescribeImages and
# should be set explicitly in configuration after import.
terraform import awsappstream_imported_image.example "example-image"
