# Example: Chronicle AZURE_BLOBSTORE_V2 Feed
# This feed type uses Google Cloud Storage Transfer Service for improved ingestion

# Example 1: Using Shared Key authentication
resource "chronicle_feed_azure_blobstore_v2" "example_shared_key" {
  display_name = "Azure Blob Storage V2 Feed - Shared Key"
  log_type     = "AZURE_AD"
  enabled      = true

  # Optional: Namespace for asset correlation
  namespace = "azure-prod"

  # Optional: Labels for feed categorization
  labels = {
    environment = "production"
    team        = "security"
  }

  details {
    azure_uri = "https://myaccount.blob.core.windows.net/logs/"

    source_delete_options = "ON_SUCCESS"  # or "NEVER"
    max_lookback_days     = 180            # Default: 180

    authentication {
      shared_key = var.azure_shared_key
      # Note: Use either shared_key OR sas_token, not both
    }
  }
}

# Example 2: Using SAS Token authentication
resource "chronicle_feed_azure_blobstore_v2" "example_sas_token" {
  display_name = "Azure Blob Storage V2 Feed - SAS Token"
  log_type     = "AZURE_AD"
  enabled      = true

  namespace = "azure-prod"

  details {
    azure_uri = "https://myaccount.blob.core.windows.net/logs/"

    source_delete_options = "NEVER"
    max_lookback_days     = 90

    authentication {
      sas_token = var.azure_sas_token
      # Note: Use either shared_key OR sas_token, not both
    }
  }
}

# Variables for sensitive data
variable "azure_shared_key" {
  description = "Azure Storage Account shared access key"
  type        = string
  sensitive   = true
}

variable "azure_sas_token" {
  description = "Azure Storage SAS (Shared Access Signature) token"
  type        = string
  sensitive   = true
}
