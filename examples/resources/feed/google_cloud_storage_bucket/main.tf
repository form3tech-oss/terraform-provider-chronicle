resource "chronicle_google_cloud_storage_bucket" "bucket" {
  display_name = "mygcsfeed"
  log_type     = "GITHUB"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    bucket_uri                = "gs://test/"
    bucket_source_type        = "FOLDERS_RECURSIVE"
    source_delete_options     = "SOURCE_DELETION_NEVER"
  }
}
