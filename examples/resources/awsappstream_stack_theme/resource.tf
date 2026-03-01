# minimal stack theme
resource "awsappstream_stack_theme" "example" {
  stack_name    = "example-stack"
  title_text    = "Example AppStream"
  theme_styling = "BLUE"

  organization_logo_s3_location = {
    s3_bucket = "my-appstream-assets"
    s3_key    = "branding/logo.png"
  }

  favicon_s3_location = {
    s3_bucket = "my-appstream-assets"
    s3_key    = "branding/favicon.ico"
  }
}

