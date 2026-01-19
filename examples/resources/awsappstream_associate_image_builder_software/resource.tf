resource "awsappstream_associate_image_builder_software" "example" {
  image_builder_arn = awsappstream_image_builder.example.arn

  software_names = [
    "Microsoft_Project_2024_Standard_64Bit",
    "Microsoft_Office_2024_LTSC_Professional_Plus_64Bit",
  ]

  deploy = true
}
