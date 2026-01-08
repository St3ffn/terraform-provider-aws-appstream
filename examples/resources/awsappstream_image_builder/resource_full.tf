# full image builder
resource "awsappstream_image_builder" "example" {
  name          = "example-image-builder"
  display_name  = "Example Image Builder"
  description   = "Image builder used to create custom AppStream images."
  instance_type = "stream.standard.large"

  # Use image ARN (preferred for long-term stability)
  image_arn = "arn:aws:appstream:eu-west-1:123456789012:image/example-image"

  enable_default_internet_access = true
  appstream_agent_version        = "LATEST"

  iam_role_arn = "arn:aws:iam::123456789012:role/AppStreamImageBuilderRole"

  vpc_config {
    subnet_ids = [
      "subnet-0abc123def4567890"
    ]

    security_group_ids = [
      "sg-0123456789abcdef0"
    ]
  }

  domain_join_info {
    directory_name                         = "corp.example.com"
    organizational_unit_distinguished_name = "OU=AppStream,DC=corp,DC=example,DC=com"
  }

  access_endpoints = [
    {
      endpoint_type = "STREAMING"
      vpce_id       = "vpce-0abc123def4567890"
    }
  ]

  root_volume_config {
    volume_size_in_gb = 300
  }

  tags = {
    Environment = "dev"
    Project     = "appstream"
    Owner       = "platform-team"
  }
}
