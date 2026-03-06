# full imported image
resource "awsappstream_imported_image" "example" {
  name                   = "example-imported-image"
  iam_role_arn           = "arn:aws:iam::123456789012:role/AppStreamImportImageRole"
  source_ami_id          = "ami-0abc1234def567890"
  agent_software_version = "ALWAYS_LATEST"
  display_name           = "Example Imported Image"
  description            = "Imported AppStream image managed by Terraform."

  runtime_validation_config = {
    intended_instance_type = "stream.standard.large"
  }

  app_catalog_config = [
    {
      name               = "example-app"
      absolute_app_path  = "C:\\Program Files\\ExampleApp\\example.exe"
      display_name       = "Example Application"
      launch_parameters  = "--mode production"
      working_directory  = "C:\\Program Files\\ExampleApp"
      absolute_icon_path = "C:\\Program Files\\ExampleApp\\example.png"
    },
    {
      name                   = "example-app-prewarm"
      absolute_app_path      = "C:\\Program Files\\ExampleApp\\example.exe"
      absolute_manifest_path = "C:\\ProgramData\\AppStream\\example-prewarm.txt"
    }
  ]

  tags = {
    Environment = "dev"
    Project     = "appstream"
    Owner       = "platform-team"
  }
}
