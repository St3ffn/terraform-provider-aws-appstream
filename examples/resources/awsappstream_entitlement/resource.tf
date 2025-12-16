resource "awsappstream_entitlement" "example" {
  stack_name     = "example-stack"
  name           = "example-name"
  description    = "Example entitlement managed by Terraform."
  app_visibility = "ALL"

  attributes = [
    {
      name  = "roles"
      value = "admin"
    }
  ]
}
