# minimal copied image
resource "awsappstream_copied_image" "example" {
  name               = "example-copied-image"
  description        = "Copied image for development"
  source_image_name  = "example-source-image"
  destination_region = "us-east-1"

  tags = {
    Environment = "dev"
    Project     = "appstream"
  }
}
