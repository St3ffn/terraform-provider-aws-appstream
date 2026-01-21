# full sessions
data "awsappstream_sessions" "example" {
  stack_name = "example-stack"
  fleet_name = "example-fleet"

  user_id             = "example-user-id"
  authentication_type = "USERPOOL"
  instance_id         = "example-instance-id"
}
