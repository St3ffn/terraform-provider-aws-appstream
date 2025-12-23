# always on fleet - single session
resource "awsappstream_fleet" "always_on_single_session" {
  name       = "always-on-single-session"
  fleet_type = "ALWAYS_ON"

  image_name    = "example-appstream-image"
  instance_type = "stream.standard.medium"

  compute_capacity {
    desired_instances = 2
  }

  display_name = "Always-On Single-Session Fleet"
}

# always on fleet - multi session
resource "awsappstream_fleet" "always_on_multi_session" {
  name       = "always-on-multi-session"
  fleet_type = "ALWAYS_ON"

  image_name    = "example-appstream-image"
  instance_type = "stream.standard.large"

  compute_capacity {
    desired_sessions = 30
  }

  max_sessions_per_instance = 5

  display_name = "Always-On Multi-Session Fleet"
}
