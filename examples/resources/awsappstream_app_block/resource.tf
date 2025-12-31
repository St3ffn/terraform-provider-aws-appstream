# minimal app block
resource "awsappstream_app_block" "example" {
  name = "example-app-block"

  source_s3_location {
    s3_bucket = "my-appstream-assets"
    s3_key    = "app-blocks/example/source.zip"
  }

  setup_script_details {
    script_s3_location {
      s3_bucket = "my-appstream-assets"
      s3_key    = "app-blocks/example/setup.ps1"
    }

    executable_path       = "powershell.exe"
    executable_parameters = "-File setup.ps1"
    timeout_in_seconds    = 300
  }
}
