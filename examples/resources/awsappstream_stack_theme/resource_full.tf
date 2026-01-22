# full stack theme
resource "awsappstream_stack_theme" "example" {
  stack_name    = "example-stack"
  title_text    = "Example AppStream Portal"
  theme_styling = "LIGHT_BLUE"

  organization_logo_s3_location {
    s3_bucket = "my-appstream-assets"
    s3_key    = "branding/logo.png"
  }

  favicon_s3_location {
    s3_bucket = "my-appstream-assets"
    s3_key    = "branding/favicon.ico"
  }

  footer_links = [
    {
      display_name    = "IT Support"
      footer_link_url = "https://support.example.com"
    },
    {
      display_name    = "Documentation"
      footer_link_url = "https://docs.example.com/appstream"
    },
    {
      display_name    = "Privacy Policy"
      footer_link_url = "https://example.com/privacy"
    }
  ]
}

