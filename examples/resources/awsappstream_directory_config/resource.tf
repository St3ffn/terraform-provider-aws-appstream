# minimal directory config
resource "awsappstream_directory_config" "example" {
  directory_name = "corp.example.com"

  organizational_unit_distinguished_names = [
    "OU=AppStream,DC=corp,DC=example,DC=com"
  ]
}
