resource "awsappstream_associate_app_block_builder_app_block" "example" {
  app_block_builder_name = "example-app-block-builder"
  app_block_arn          = "arn:aws:appstream:eu-central-1:123456789012:app-block/example-app-block"
}
