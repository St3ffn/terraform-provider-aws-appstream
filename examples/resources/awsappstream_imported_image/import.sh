# Import uses: name|iam_role_arn|source_ami_id.
# Note: create-time-only attributes (agent_software_version, runtime_validation_config, app_catalog_config)
# are not returned by DescribeImages and should be set explicitly in configuration after import.
terraform import awsappstream_imported_image.example "example-image|arn:aws:iam::123456789012:role/appstream-import|ami-0abc1234def567890"
