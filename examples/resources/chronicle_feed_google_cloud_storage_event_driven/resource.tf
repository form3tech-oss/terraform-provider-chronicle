# Example: Chronicle GOOGLE_CLOUD_STORAGE_EVENT_DRIVEN Feed
# This feed type uses Google Cloud Storage Transfer Service with Pub/Sub notifications
# for event-driven, push-based ingestion with reduced latency

resource "chronicle_feed_google_cloud_storage_event_driven" "example" {
  display_name = "GCS Event-Driven Feed Example"
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
    bucket_uri = "gs://my-security-logs/cloudaudit/"

    # Pub/Sub subscription created for Cloud Storage notifications
    # Format: projects/PROJECT_ID/subscriptions/SUBSCRIPTION_NAME
    pubsub_subscription = "projects/my-project/subscriptions/my-gcs-notifications"

    source_delete_options = "ON_SUCCESS"  # or "NEVER"
    max_lookback_days     = 180            # Default: 180, Max: 180
  }
}

# Prerequisites for this feed type:
# 1. Create a Pub/Sub topic for Cloud Storage notifications
# 2. Configure Cloud Storage bucket to send notifications to the topic
# 3. Create a Pub/Sub subscription on the topic
# 4. Grant the Google Security Operations service account:
#    - Storage Object Viewer role on the bucket
#    - Pub/Sub Subscriber role on the subscription
#
# Example setup commands:
#
# # Get the Google Security Operations service account
# # Use Chronicle Feed Management API: fetchFeedServiceAccount
#
# # Create Pub/Sub topic
# gcloud pubsub topics create my-gcs-notifications
#
# # Configure bucket notifications
# gsutil notification create -t my-gcs-notifications -f json gs://my-security-logs
#
# # Create subscription
# gcloud pubsub subscriptions create my-gcs-notifications \
#   --topic=my-gcs-notifications \
#   --ack-deadline=600
#
# # Grant permissions
# gsutil iam ch serviceAccount:SERVICE_ACCOUNT@PROJECT.iam.gserviceaccount.com:objectViewer \
#   gs://my-security-logs
#
# gcloud pubsub subscriptions add-iam-policy-binding my-gcs-notifications \
#   --member=serviceAccount:SERVICE_ACCOUNT@PROJECT.iam.gserviceaccount.com \
#   --role=roles/pubsub.subscriber
