# export image tasks filtered by image ARN
data "awsappstream_export_image_tasks" "example" {
  filters = [{
    name   = "image-arn"
    values = ["my-image"]
  }]
}
