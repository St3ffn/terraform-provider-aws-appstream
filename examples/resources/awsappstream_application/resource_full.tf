# full resource
resource "awsappstream_application" "example" {
  name         = "example-app"
  display_name = "Example Application"
  description  = "Example AppStream application managed by Terraform"

  icon_s3_location {
    s3_bucket = "my-appstream-assets"
    s3_key    = "icons/example.png"
  }

  launch_path       = "C:\\Program Files\\ExampleApp\\example.exe"
  working_directory = "C:\\Program Files\\ExampleApp"
  launch_parameters = "--mode production"

  platforms = [
    "WINDOWS_SERVER_2019",
    "WINDOWS_SERVER_2022",
  ]

  instance_families = [
    "stream.standard",
    "stream.compute",
  ]

  app_block_arn = "arn:aws:appstream:eu-west-1:123456789012:app-block/example-block"

  tags = {
    Name        = "example-app"
    Environment = "production"
    Owner       = "platform-team"
  }
}
