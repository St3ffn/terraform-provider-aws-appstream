resource "awsappstream_image_permission" "example" {
  name              = "example-image"
  shared_account_id = "123456789012"

  image_permissions {
    allow_fleet         = true
    allow_image_builder = true
  }
}
