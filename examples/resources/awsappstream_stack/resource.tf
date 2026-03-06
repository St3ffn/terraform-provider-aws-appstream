# minimal stack
resource "awsappstream_stack" "example" {
  name = "example-stack"

  tags = {
    Environment = "dev"
    Project     = "appstream"
  }
}
