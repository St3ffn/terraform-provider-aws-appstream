# full user
resource "awsappstream_user" "full" {
  authentication_type = "USERPOOL"
  user_name           = "example@example.com"

  first_name = "Example"
  last_name  = "User"

  # Write-only, applies only at creation
  message_action = "SUPPRESS"

  # Explicitly manage enabled state
  enabled = false
}

