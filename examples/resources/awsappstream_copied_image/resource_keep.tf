# copied image protected from destroy
resource "awsappstream_copied_image" "keep" {
  destination_image_name        = "example-copied-image-keep"
  destination_image_description = "Copied image retained for rollback"
  destination_region            = "us-east-1"
  source_image_name             = "example-source-image"

  # Prevent accidental image deletion via `terraform destroy` or replacement plans.
  lifecycle {
    prevent_destroy = true
  }

  tags = {
    Environment = "dev"
    Project     = "appstream"
    Retention   = "keep"
  }
}
