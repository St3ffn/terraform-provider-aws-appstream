resource "awsappstream_associate_user_stack" "example" {
  stack_name          = "example-stack"
  user_name           = "example@example.com"
  authentication_type = "USERPOOL"

  send_email_notification = true
}
