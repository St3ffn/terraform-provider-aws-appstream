data "awsappstream_user" "example" {
  authentication_type = "USERPOOL"
  user_name           = "example@example.com"
}
