# minimal app block builder
resource "awsappstream_app_block_builder" "example" {
  name          = "example-app-block-builder"
  instance_type = "stream.standard.large"
  platform      = "WINDOWS_SERVER_2019"

  vpc_config {
    # subnets must be in different Availability Zones
    subnet_ids = [
      "subnet-0abc123def4567890",
      "subnet-1def9876543210987"
    ]
  }
}
