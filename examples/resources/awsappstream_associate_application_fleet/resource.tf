resource "awsappstream_associate_application_fleet" "example" {
  fleet_name      = "example-fleet"
  application_arn = "arn:aws:appstream:eu-west-1:123456789012:application/example-app"
}
