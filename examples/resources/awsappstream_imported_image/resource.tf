# minimal imported image
resource "awsappstream_imported_image" "example" {
  name          = "example-imported-image"
  iam_role_arn  = "arn:aws:iam::123456789012:role/AppStreamImportImageRole"
  source_ami_id = "ami-0abc1234def567890"

  tags = {
    Environment = "dev"
    Project     = "appstream"
  }
}
