# imported image protected from destroy
resource "awsappstream_imported_image" "keep" {
  name          = "example-imported-image-keep"
  iam_role_arn  = "arn:aws:iam::123456789012:role/AppStreamImportImageRole"
  source_ami_id = "ami-0abc1234def567890"

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
