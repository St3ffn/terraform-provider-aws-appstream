# appstream image by arn
data "awsappstream_image" "by_arn" {
  arn = "arn:aws:appstream:us-east-1:123456789012:image/my-image"
}

# appstream image by name
data "awsappstream_image" "by_name" {
  name = "my-appstream-image"
}

# appstream image by name_regex
data "awsappstream_image" "by_regex" {
  name_regex = "^my-appstream-image-.*"
}

# most recent appstream image by name_regex
data "awsappstream_image" "latest" {
  name_regex  = "^my-appstream-image-"
  most_recent = true
}

# most recent appstream image by name_regex with visibility PRIVATE
data "awsappstream_image" "latest_private" {
  name_regex  = "^golden-image-"
  visibility  = "PRIVATE"
  most_recent = true
}
