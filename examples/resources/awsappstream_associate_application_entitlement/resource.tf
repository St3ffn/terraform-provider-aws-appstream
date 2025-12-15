resource "awsappstream_associate_application_entitlement" "example" {
  stack_name             = "example-stack"
  entitlement_name       = "example-entitlement"
  application_identifier = "example-application"
}
