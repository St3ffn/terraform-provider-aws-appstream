# minimal image builder
resource "awsappstream_image_builder" "example" {
  name          = "example-image-builder"
  instance_type = "stream.standard.medium"

  # Exactly one of image_name or image_arn is required
  image_name = "AppStream-Windows-Server-2022"
}
