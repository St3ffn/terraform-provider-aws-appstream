#!/usr/bin/env bash
# Copyright (c) St3ffn
# SPDX-License-Identifier: MPL-2.0

set -euo pipefail

REGION="${AWS_REGION:-${AWS_DEFAULT_REGION:-}}"

if [[ -z "$REGION" ]]; then
  echo "ERROR: AWS_REGION or AWS_DEFAULT_REGION must be set"
  exit 1
fi

ACCOUNT_ID="$(aws sts get-caller-identity --query Account --output text)"

BUCKET_NAME="appstream-acc-test-bucket-${ACCOUNT_ID}-${REGION}"
BUCKET_URI="s3://$BUCKET_NAME"

APP_BLOCK_VHD_FILENAME="appblock.vhd"
APP_BLOCK_SETUP_SCRIPT_FILENAME="app_block_setup.sh"
APP_BLOCK_POST_SETUP_SCRIPT_FILENAME="app_block_post_setup.sh"

APPLICATION_ICON_FILENAME="application_icon.png"

echo "Setting up AppStream acceptance test prerequisites"
echo "Account: $ACCOUNT_ID"
echo "Region:  $REGION"
echo "Bucket:  $BUCKET_NAME"
echo

# Create bucket if missing
if aws s3api head-bucket --bucket "$BUCKET_NAME" >/dev/null 2>&1; then
  echo "✓ Bucket exists"
else
  echo "→ Creating bucket..."
  if [[ "$REGION" == "us-east-1" ]]; then
    aws s3api create-bucket \
      --bucket "$BUCKET_NAME"
  else
    aws s3api create-bucket \
      --bucket "$BUCKET_NAME" \
      --create-bucket-configuration LocationConstraint="$REGION"
  fi
  echo "→ Created bucket"
fi

# Create bucket policy if missing
if aws s3api get-bucket-policy --bucket "$BUCKET_NAME" >/dev/null 2>&1; then
  echo "✓ Bucket policy already exists"
else
  echo "→ Applying bucket policy..."
  aws s3api put-bucket-policy \
    --bucket "$BUCKET_NAME" \
    --policy "$(cat <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAppStreamListBucket",
      "Effect": "Allow",
      "Principal": {
        "Service": "appstream.amazonaws.com"
      },
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::$BUCKET_NAME",
      "Condition": {
        "StringEquals": {
          "aws:SourceAccount": "$ACCOUNT_ID"
        }
      }
    },
    {
      "Sid": "AllowAppStreamRead",
      "Effect": "Allow",
      "Principal": {
        "Service": "appstream.amazonaws.com"
      },
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::$BUCKET_NAME/*",
      "Condition": {
        "StringEquals": {
          "aws:SourceAccount": "$ACCOUNT_ID"
        }
      }
    },
    {
      "Sid": "AllowAppStreamBucketOwnershipControls",
      "Effect": "Allow",
      "Principal": {
        "Service": "appstream.amazonaws.com"
      },
      "Action": "s3:GetBucketOwnershipControls",
      "Resource": "arn:aws:s3:::$BUCKET_NAME",
      "Condition": {
        "StringEquals": {
          "aws:SourceAccount": "$ACCOUNT_ID"
        }
      }
    }
  ]
}
EOF
)"
  echo "→ Bucket policy applied"
fi

# Temporary Directory
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

# Create app block vhd if missing
APP_BLOCK_VHD_FILE="${TMP_DIR}/${APP_BLOCK_VHD_FILENAME}"

if aws s3api head-object --bucket "$BUCKET_NAME" --key "$APP_BLOCK_VHD_FILENAME" >/dev/null 2>&1; then
  echo "✓ $APP_BLOCK_VHD_FILENAME already exists in bucket"
else
  echo "→ Creating $APP_BLOCK_VHD_FILENAME..."
  dd if=/dev/zero of="$APP_BLOCK_VHD_FILE" bs=1M count=1 status=none

  echo "→ Uploading $APP_BLOCK_VHD_FILENAME..."
  aws s3 cp "$APP_BLOCK_VHD_FILE" "$BUCKET_URI/$APP_BLOCK_VHD_FILENAME"

  echo "✓ $APP_BLOCK_VHD_FILENAME uploaded"
fi

# Create app block setup script if missing
APP_BLOCK_SETUP_SCRIPT_FILE="${TMP_DIR}/${APP_BLOCK_SETUP_SCRIPT_FILENAME}"

if aws s3api head-object --bucket "$BUCKET_NAME" --key "$APP_BLOCK_SETUP_SCRIPT_FILENAME" >/dev/null 2>&1; then
  echo "✓ $APP_BLOCK_SETUP_SCRIPT_FILENAME already exists in bucket"
else
  echo "→ Creating $APP_BLOCK_SETUP_SCRIPT_FILENAME..."

  cat <<'EOF' > "$APP_BLOCK_SETUP_SCRIPT_FILE"
#!/bin/sh
echo "app block setup script executed"
EOF

  echo "→ Uploading $APP_BLOCK_SETUP_SCRIPT_FILENAME..."
  aws s3 cp "$APP_BLOCK_SETUP_SCRIPT_FILE" "$BUCKET_URI/$APP_BLOCK_SETUP_SCRIPT_FILENAME"

  echo "✓ $APP_BLOCK_SETUP_SCRIPT_FILENAME uploaded"
fi

# Create app block setup post script if missing
APP_BLOCK_POST_SETUP_SCRIPT_FILE="${TMP_DIR}/${APP_BLOCK_POST_SETUP_SCRIPT_FILENAME}"

if aws s3api head-object --bucket "$BUCKET_NAME" --key "$APP_BLOCK_POST_SETUP_SCRIPT_FILENAME" >/dev/null 2>&1; then
  echo "✓ $APP_BLOCK_POST_SETUP_SCRIPT_FILENAME already exists in bucket"
else
  echo "→ Creating $APP_BLOCK_POST_SETUP_SCRIPT_FILENAME..."

  cat <<'EOF' > "$APP_BLOCK_POST_SETUP_SCRIPT_FILE"
#!/bin/sh
echo "app block post setup script executed"
EOF

  echo "→ Uploading $APP_BLOCK_POST_SETUP_SCRIPT_FILENAME..."
  aws s3 cp "$APP_BLOCK_POST_SETUP_SCRIPT_FILE" "$BUCKET_URI/$APP_BLOCK_POST_SETUP_SCRIPT_FILENAME"

  echo "✓ $APP_BLOCK_POST_SETUP_SCRIPT_FILENAME uploaded"
fi

# Create application icon if missing
APPLICATION_ICON_FILE="${TMP_DIR}/${APPLICATION_ICON_FILENAME}"

if aws s3api head-object --bucket "$BUCKET_NAME" --key "$APPLICATION_ICON_FILENAME" >/dev/null 2>&1; then
  echo "✓ $APPLICATION_ICON_FILENAME already exists in bucket"
else
  echo "→ Creating $APPLICATION_ICON_FILENAME..."

  cat <<'EOF' | base64 --decode > "$APPLICATION_ICON_FILE"
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO7+ZxkAAAAASUVORK5CYII=
EOF

  echo "→ Uploading $APPLICATION_ICON_FILENAME..."
  aws s3 cp "$APPLICATION_ICON_FILE" "$BUCKET_URI/$APPLICATION_ICON_FILENAME"

  echo "✓ $APPLICATION_ICON_FILENAME uploaded"
fi
