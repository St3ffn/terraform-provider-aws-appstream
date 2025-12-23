# on demand fleet - single session
resource "awsappstream_fleet" "on_demand_single_session" {
  name       = "on-demand-single-session"
  fleet_type = "ON_DEMAND"

  image_name    = "example-appstream-image"
  instance_type = "stream.standard.medium"

  compute_capacity {
    desired_instances = 2
  }

  display_name = "On-Demand Single-Session Fleet"
}

# on demand fleet - multi session
resource "awsappstream_fleet" "on_demand_multi_session" {
  name       = "on-demand-multi-session"
  fleet_type = "ON_DEMAND"

  image_name    = "example-appstream-image"
  instance_type = "stream.standard.large"

  compute_capacity {
    desired_sessions = 20
  }

  max_sessions_per_instance = 4

  display_name = "On-Demand Multi-Session Fleet"
}
