# full provider configuration
provider "awsappstream" {
  profile            = "aws-test-account-profile"
  region             = "eu-central-1" # required
  retry_max_attempts = 10
  retry_max_backoff  = 30
}
