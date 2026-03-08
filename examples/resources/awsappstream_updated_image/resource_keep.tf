# updated image protected from destroy
# This operation can trigger AWS-managed image-builder processing and may take a long time.
resource "awsappstream_updated_image" "keep" {
  existing_image_name = "example-existing-image"

  new_image_name         = "example-new-image"
  new_image_description  = "my new image with installed updates"
  new_image_display_name = "example-new-updated-image"

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
