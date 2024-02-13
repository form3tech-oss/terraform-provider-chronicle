resource "chronicle_feed_amazon_s3" "s3" {
  display_name = "mys3feed"
  log_type     = "GITHUB"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    s3_uri                = "s3://s3-bucket/"
    s3_source_type        = "FOLDERS_RECURSIVE"
    source_delete_options = "SOURCE_DELETION_NEVER"
    authentication {
      region            = "EU_WEST_1"
      access_key_id     = "XXXXX"
      secret_access_key = "XXXX"
    }
  }

}
