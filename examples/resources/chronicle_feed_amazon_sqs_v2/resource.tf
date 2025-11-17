# Example: Chronicle AMAZON_SQS_V2 Feed
# This feed type uses Google Cloud Storage Transfer Service for improved ingestion

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
    s3_uri         = "s3://my-security-logs/cloudtrail/"
    region         = "us-east-1"
    account_number = "123456789012"
    queue_name     = "my-cloudtrail-queue"

    source_delete_options = "ON_SUCCESS"  # or "NEVER"
    max_lookback_days     = 180            # Default: 180

    authentication {
      sqs_access_key_id     = var.sqs_access_key_id
      sqs_secret_access_key = var.sqs_secret_access_key

      # Optional: Specify separate credentials for S3 bucket access
      # s3_authentication {
      #   access_key_id     = var.s3_access_key_id
      #   secret_access_key = var.s3_secret_access_key
      # }
    }
  }
}

# Variables for sensitive data
variable "sqs_access_key_id" {
  description = "AWS access key ID for SQS queue"
  type        = string
  sensitive   = true
}

variable "sqs_secret_access_key" {
  description = "AWS secret access key for SQS queue"
  type        = string
  sensitive   = true
}

# Optional variables for separate S3 credentials
# variable "s3_access_key_id" {
#   description = "AWS access key ID for S3 bucket"
#   type        = string
#   sensitive   = true
# }

# variable "s3_secret_access_key" {
#   description = "AWS secret access key for S3 bucket"
#   type        = string
#   sensitive   = true
# }
