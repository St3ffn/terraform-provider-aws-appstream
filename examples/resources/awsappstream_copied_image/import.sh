# Import uses destination image name and destination region.
# Note: create-time-only attribute source_image_name should be set explicitly in configuration after import.
terraform import awsappstream_copied_image.example "example-copied-image|us-east-1"
