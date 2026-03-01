# full provider configuration
provider "awsappstream" {
  profile = "appstream-admin"
  region  = "eu-central-1"

  retry_mode         = "adaptive"
  retry_max_attempts = 10
  retry_max_backoff  = 30

  default_tags {
    tags = {
      maintainer  = "st3ffn"
      environment = "prod"
      managed_by  = "terraform"
    }
  }
}
