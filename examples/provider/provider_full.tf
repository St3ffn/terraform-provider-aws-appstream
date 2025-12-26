# full provider configuration
provider "awsappstream" {
  profile = "appstream-admin"
  region  = "eu-central-1"

  retry_mode         = "adaptive"
  retry_max_attempts = 10
  retry_max_backoff  = 30

  default_tags {
    tags = {
      environment = "prod"
      team        = "platform"
      managed_by  = "terraform"
    }
  }
}
