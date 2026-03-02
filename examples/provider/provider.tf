terraform {
  required_version = ">= 1.2"

  required_providers {
    awsappstream = {
      source  = "st3ffn/aws-appstream"
      version = "~> 2.0"
    }
  }
}

# minimal provider configuration
provider "awsappstream" {
  region = "eu-central-1" # required
}
