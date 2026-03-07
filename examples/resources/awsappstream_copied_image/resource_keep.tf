# copied image protected from destroy
resource "awsappstream_copied_image" "keep" {
  name               = "example-copied-image-keep"
  description        = "Copied image retained for rollback"
  source_image_name  = "example-source-image"
  destination_region = "us-east-1"

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
