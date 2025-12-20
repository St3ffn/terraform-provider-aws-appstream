# full resource
resource "awsappstream_stack" "example" {
  name         = "example-stack"
  display_name = "Example AppStream Stack"
  description  = "Example AppStream stack managed by Terraform."

  redirect_url = "https://example.com/logout"
  feedback_url = "https://example.com/feedback"

  storage_connectors = [
    {
      connector_type = "HOMEFOLDERS"
    },
    {
      connector_type                = "ONE_DRIVE"
      domains                       = ["example.com"]
      domains_require_admin_consent = ["example.com"]
    }
  ]

  user_settings = [
    {
      action     = "FILE_UPLOAD"
      permission = "ENABLED"
    },
    {
      action         = "CLIPBOARD_COPY_TO_LOCAL_DEVICE"
      permission     = "ENABLED"
      maximum_length = 1048576
    }
  ]

  application_settings = {
    enabled        = true
    settings_group = "example-settings"
  }

  access_endpoints = [
    {
      endpoint_type = "STREAMING"
      vpce_id       = "vpce-0abc123def4567890"
    }
  ]

  embed_host_domains = [
    "apps.example.com",
    "portal.example.com"
  ]

  streaming_experience_settings = {
    preferred_protocol = "TCP"
  }

  tags = {
    Environment = "dev"
    Project     = "appstream"
    Owner       = "platform-team"
  }
}
