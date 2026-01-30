# Example: Chronicle AMAZON_SQS_V2 Feed
# This feed type uses Google Cloud Storage Transfer Service for improved ingestion
# SQS V2 provides push-based ingestion with reduced latency compared to S3 V2

resource "chronicle_feed_amazon_sqs_v2" "example" {
  display_name = "AWS SQS V2 Feed Example"
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
    # Full ARN of the SQS queue
    # Format: arn:aws:sqs:region:account_id:queue_name
    queue  = "arn:aws:sqs:us-east-1:123456789012:my-cloudtrail-queue"
    s3_uri = "s3://my-security-logs/cloudtrail/"

    source_delete_options = "ON_SUCCESS"  # or "NEVER"
    max_lookback_days     = 180            # Default: 180, Max: 180

    # Authentication using access keys
    # Both SQS queue and S3 bucket use the same credentials
    authentication {
      access_key_id     = var.aws_access_key_id
      secret_access_key = var.aws_secret_access_key
    }
  }
}

# Example using IAM role ARN (Federated authentication)
resource "chronicle_feed_amazon_sqs_v2" "example_iam" {
  display_name = "AWS SQS V2 Feed with IAM Role"
  log_type     = "AWS_CLOUDTRAIL"
  enabled      = true
  namespace    = "aws-prod"

  details {
    queue  = "arn:aws:sqs:us-east-1:123456789012:my-cloudtrail-queue"
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
  description = "AWS access key ID for both SQS queue and S3 bucket"
  type        = string
  sensitive   = true
}

variable "aws_secret_access_key" {
  description = "AWS secret access key for both SQS queue and S3 bucket"
  type        = string
  sensitive   = true
}
