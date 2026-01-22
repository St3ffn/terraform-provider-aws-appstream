# full app block builder
resource "awsappstream_app_block_builder" "example" {
  name          = "example-app-block-builder"
  display_name  = "Example App Block Builder"
  description   = "App block builder used to package applications or app blocks"
  instance_type = "stream.standard.large"
  platform      = "WINDOWS_SERVER_2019"

  enable_default_internet_access = true

  iam_role_arn = "arn:aws:iam::123456789012:role/AppStreamAppBlockBuilderRole"

  vpc_config {
    # subnets must be in different Availability Zones
    subnet_ids = [
      "subnet-0abc123def4567890",
      "subnet-1def9876543210987"
    ]

    security_group_ids = [
      "sg-0123456789abcdef0"
    ]
  }

  access_endpoints = [
    {
      endpoint_type = "STREAMING"
      vpce_id       = "vpce-0abc123def4567890"
    }
  ]

  tags = {
    Environment = "dev"
    Project     = "appstream"
    Owner       = "platform-team"
  }
}
