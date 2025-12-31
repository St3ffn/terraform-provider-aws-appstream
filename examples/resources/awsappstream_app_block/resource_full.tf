# full app block - CUSTOM packaging
resource "awsappstream_app_block" "example" {
  name         = "example-app-block"
  display_name = "Example App Block"
  description  = "Reusable AppStream app block managed by Terraform"

  packaging_type = "CUSTOM"

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
    executable_parameters = "-ExecutionPolicy Bypass -File setup.ps1"
    timeout_in_seconds    = 600
  }

  tags = {
    Name        = "example-app-block"
    Environment = "production"
    Owner       = "platform-team"
  }
}

# full app block - APPSTREAM2 packaging
resource "awsappstream_app_block" "example" {
  name           = "example-app-block"
  packaging_type = "APPSTREAM2"

  source_s3_location {
    s3_bucket = "my-appstream-assets"
    # s3_key optional when AppStream builds a new package
  }

  post_setup_script_details {
    script_s3_location {
      s3_bucket = "my-appstream-assets"
      s3_key    = "app-blocks/example/post-setup.ps1"
    }

    executable_path       = "powershell.exe"
    executable_parameters = "-File post-setup.ps1"
    timeout_in_seconds    = 300
  }
}
