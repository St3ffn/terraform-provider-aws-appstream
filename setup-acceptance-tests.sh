#!/usr/bin/env bash
# Copyright (c) St3ffn
# SPDX-License-Identifier: MPL-2.0

set -euo pipefail

REGION="${AWS_REGION:-eu-central-1}"
ACCOUNT_ID="$(aws sts get-caller-identity --query Account --output text)"

BUCKET_NAME="appstream-acc-test-bucket"
BUCKET_URI="s3://$BUCKET_NAME"

APP_BLOCK_VHD_FILENAME="appblock.vhd"
APP_BLOCK_SETUP_SCRIPT_FILENAME="app_block_setup.sh"
APP_BLOCK_POST_SETUP_SCRIPT_FILENAME="app_block_post_setup.sh"

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
  aws s3api create-bucket \
    --bucket "$BUCKET_NAME" \
    --region "$REGION" \
    --create-bucket-configuration LocationConstraint="$REGION"
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
