# elastic fleet
resource "awsappstream_fleet" "elastic" {
  name       = "elastic-fleet"
  fleet_type = "ELASTIC"

  instance_type = "stream.standard.large"
  platform      = "WINDOWS_SERVER_2019"

  max_concurrent_sessions = 100

  vpc_config = {
    subnet_ids = [
      "subnet-0123456789abcdef0",
      "subnet-0fedcba9876543210",
    ]
  }

  desired_state = "RUNNING"

  display_name = "Elastic Fleet"

  tags = {
    Environment = "dev"
    Project     = "appstream"
  }
}
