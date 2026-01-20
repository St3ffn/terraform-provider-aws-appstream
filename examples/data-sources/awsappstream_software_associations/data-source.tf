# appstream image
data "awsappstream_software_associations" "example_image" {
  associated_resource = "arn:aws:appstream:eu-central-1:123456789012:image/example-image"
}

# appstream image-builder
data "awsappstream_software_associations" "example_image_builder" {
  associated_resource = "arn:aws:appstream:eu-central-1:123456789012:image-builder/example-image-builder"
}
