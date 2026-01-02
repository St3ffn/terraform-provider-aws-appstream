# minimal user
resource "awsappstream_user" "minimal" {
  authentication_type = "USERPOOL"
  user_name           = "example@example.com"
}

