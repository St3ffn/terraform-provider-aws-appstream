# full directory config
resource "awsappstream_directory_config" "example" {
  directory_name = "corp.example.com"

  organizational_unit_distinguished_names = [
    "OU=AppStream,DC=corp,DC=example,DC=com",
    "OU=ImageBuilders,DC=corp,DC=example,DC=com"
  ]

  service_account_credentials {
    account_name     = "svc-appstream"
    account_password = var.appstream_service_account_password
  }

  certificate_based_auth_properties {
    status = "ENABLED"

    certificate_authority_arn = "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abcd1234-5678-90ab-cdef-EXAMPLE11111"
  }
}
