# minimal sessions
data "awsappstream_sessions" "example" {
  stack_name = "example-stack"
  fleet_name = "example-fleet"
}
