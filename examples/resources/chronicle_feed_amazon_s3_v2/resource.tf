# Example: Chronicle AMAZON_S3_V2 Feed
# This feed type uses Google Cloud Storage Transfer Service for improved ingestion

resource "chronicle_feed_amazon_s3_v2" "example" {
  display_name = "AWS S3 V2 Feed Example"
  log_type     = "AWS_CLOUDTRAIL"
  enabled      = true

  # Optional: Namespace for asset correlation
  namespace = "aws-prod"

  # Optional: Labels for feed categorization
  labels = {
    environment = "production"
    team        = "security"
  }

  details {
    s3_uri = "s3://my-security-logs/cloudtrail/"

    source_delete_options = "ON_SUCCESS"  # or "NEVER"
    max_lookback_days     = 180            # Default: 180, Max: 180

    # Authentication using access keys
    authentication {
      access_key_id     = var.aws_access_key_id
      secret_access_key = var.aws_secret_access_key
    }
  }
}

# Example using IAM role ARN (Federated authentication)
resource "chronicle_feed_amazon_s3_v2" "example_iam" {
  display_name = "AWS S3 V2 Feed with IAM Role"
  log_type     = "AWS_CLOUDTRAIL"
  enabled      = true
  namespace    = "aws-prod"

  details {
    s3_uri = "s3://my-security-logs/cloudtrail/"

    source_delete_options = "ON_SUCCESS"
    max_lookback_days     = 90

    # Authentication using IAM role (recommended for federated access)
    authentication {
      aws_iam_role_arn = "arn:aws:iam::123456789012:role/chronicle-ingestion-role"
    }
  }
}

# Variables for sensitive data
variable "aws_access_key_id" {
  description = "AWS access key ID (20 characters)"
  type        = string
  sensitive   = true
}

variable "aws_secret_access_key" {
  description = "AWS secret access key (40 characters)"
  type        = string
  sensitive   = true
}
