# minimal copied image
resource "awsappstream_copied_image" "example" {
  destination_image_name        = "example-copied-image"
  destination_image_description = "Copied image for development"
  destination_region            = "us-east-1"
  source_image_name             = "example-source-image"

  tags = {
    Environment = "dev"
    Project     = "appstream"
  }
}
