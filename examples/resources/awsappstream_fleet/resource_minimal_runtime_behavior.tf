# minimal fleet with runtime behavior controls
resource "awsappstream_fleet" "minimal_runtime_behavior" {
  name          = "minimal-runtime-behavior-fleet"
  fleet_type    = "ON_DEMAND"
  image_name    = "example-image"
  instance_type = "stream.standard.small"

  compute_capacity = {
    desired_instances = 1
  }

  # desired_state enforces the fleet runtime state after create/update (INHERIT, RUNNING, STOPPED).
  desired_state = "STOPPED"
  # update_behavior controls update handling when AWS requires stop: auto stop/start or fail if running.
  update_behavior = "FAIL_IF_RUNNING"

  tags = {
    Environment = "dev"
    Project     = "appstream"
  }
}
