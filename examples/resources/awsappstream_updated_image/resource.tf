# minimal update image
# This operation can trigger AWS-managed image-builder processing and may take a long time.
resource "awsappstream_updated_image" "example" {
  existing_image_name = "example-existing-image"
  new_image_name      = "example-new-image"

  tags = {
    Environment = "dev"
    Project     = "appstream"
  }
}
