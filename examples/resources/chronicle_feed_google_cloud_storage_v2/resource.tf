# Example: Chronicle GOOGLE_CLOUD_STORAGE_V2 Feed
# This feed type uses Google Cloud Storage Transfer Service for improved ingestion
# Authentication is handled via the Google Security Operations service account

resource "chronicle_feed_google_cloud_storage_v2" "example" {
  display_name = "GCS V2 Feed Example"
  log_type     = "GCP_CLOUDAUDIT"
  enabled      = true

  # Optional: Namespace for asset correlation
  namespace = "gcp-prod"

  # Optional: Labels for feed categorization
  labels = {
    environment = "production"
    team        = "security"
  }

  details {
    bucket_uri = "gs://my-security-logs/audit-logs/"

    source_delete_options = "NEVER"       # or "ON_SUCCESS"
    max_lookback_days     = 180           # Default: 180
  }
}

# Note: Before setting up this feed, you must:
# 1. Get the Google Security Operations service account using the
#    fetchFeedServiceAccount method from the Feed Management API
# 2. Grant the service account access to your GCS bucket with appropriate permissions
#    (e.g., Storage Object Viewer for reading, Storage Object Admin for deletion)
