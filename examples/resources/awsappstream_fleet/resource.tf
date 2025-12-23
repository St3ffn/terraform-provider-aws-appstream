# minimal fleet
resource "awsappstream_fleet" "minimal" {
  name          = "minimal-fleet"
  fleet_type    = "ON_DEMAND"
  image_name    = "example-image"
  instance_type = "stream.standard.small"

  compute_capacity {
    desired_instances = 1
  }
}

# full fleet
resource "awsappstream_fleet" "full" {
  name       = "always-on-engineering-fleet"
  fleet_type = "ALWAYS_ON"

  image_name    = "appstream-windows-2022-engineering"
  instance_type = "stream.standard.large"

  display_name = "Engineering Always-On Fleet"
  description  = "Always-on AppStream fleet for engineering workloads"

  compute_capacity {
    desired_instances = 5
  }

  vpc_config {
    subnet_ids = [
      "subnet-0a1234567890abcd1",
      "subnet-0b1234567890abcd2",
    ]

    security_group_ids = [
      "sg-0123456789abcdef0",
    ]
  }

  max_user_duration_in_seconds       = 28800 # 8 hours
  disconnect_timeout_in_seconds      = 900   # 15 minutes
  idle_disconnect_timeout_in_seconds = 1800  # 30 minutes

  enable_default_internet_access = true

  iam_role_arn = "arn:aws:iam::123456789012:role/AppStreamFleetRole"

  stream_view = "DESKTOP"
  platform    = "WINDOWS_SERVER_2022"

  domain_join_info {
    directory_name                         = "corp.example.com"
    organizational_unit_distinguished_name = "OU=AppStream,OU=Computers,DC=corp,DC=example,DC=com"
  }

  root_volume_config {
    volume_size_in_gb = 250
  }

  usb_device_filter_strings = [
    "USB\\VID_046D&PID_C52B,*,*,*,*,*,1,0",
    "USB\\VID_0781&PID_558A,*,*,*,*,*,1,1",
  ]

  tags = {
    Environment = "production"
    Team        = "engineering"
    ManagedBy   = "terraform"
  }
}

