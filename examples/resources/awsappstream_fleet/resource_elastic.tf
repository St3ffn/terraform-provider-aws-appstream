# elastic fleet
resource "awsappstream_fleet" "elastic" {
  name       = "elastic-fleet"
  fleet_type = "ELASTIC"

  image_name    = "example-elastic-image"
  instance_type = "stream.standard.large"

  max_concurrent_sessions = 100

  vpc_config {
    subnet_ids = [
      "subnet-0123456789abcdef0",
      "subnet-0fedcba9876543210",
    ]
  }

  display_name = "Elastic Fleet"
}
