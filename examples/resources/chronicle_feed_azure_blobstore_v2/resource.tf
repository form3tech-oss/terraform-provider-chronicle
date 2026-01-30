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
    max_lookback_days     = 180            # Default: 180, Max: 180

    authentication {
      shared_key = var.azure_shared_key
      # Note: Use one of: shared_key, sas_token, or workload_identity_federation
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
      # Note: Use one of: shared_key, sas_token, or workload_identity_federation
    }
  }
}

# Example 3: Using Workload Identity Federation (recommended for production)
resource "chronicle_feed_azure_blobstore_v2" "example_federated" {
  display_name = "Azure Blob Storage V2 Feed - Federated Identity"
  log_type     = "AZURE_AD"
  enabled      = true

  namespace = "azure-prod"

  details {
    azure_uri = "https://myaccount.blob.core.windows.net/logs/"

    source_delete_options = "ON_SUCCESS"
    max_lookback_days     = 180

    authentication {
      workload_identity_federation {
        client_id = var.azure_client_id
        tenant_id = var.azure_tenant_id
      }
      # Note: Use one of: shared_key, sas_token, or workload_identity_federation
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

variable "azure_client_id" {
  description = "Application (client) ID of the registered Azure application"
  type        = string
  sensitive   = false  # Client ID is not sensitive, but can be marked sensitive if desired
}

variable "azure_tenant_id" {
  description = "Directory (tenant) ID of the registered Azure application"
  type        = string
  sensitive   = false  # Tenant ID is not sensitive, but can be marked sensitive if desired
}
