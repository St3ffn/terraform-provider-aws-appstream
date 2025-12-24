# minimal application
resource "awsappstream_application" "example" {
  name = "example-app"

  icon_s3_location {
    s3_bucket = "my-appstream-assets"
    s3_key    = "icons/example.png"
  }

  launch_path = "C:\\Program Files\\ExampleApp\\example.exe"

  platforms = [
    "WINDOWS_SERVER_2019",
  ]

  instance_families = [
    "stream.standard",
  ]

  app_block_arn = "arn:aws:appstream:eu-west-1:123456789012:app-block/example-block"
}


