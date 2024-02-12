resource "chronicle_feed_amazon_sqs" "sqs" {
  display_name = "mys3feed"
  log_type     = "GITHUB"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }

  details {
    queue = "sqs"
    region = "EU_WEST_1"
    account_number = "111"
    source_delete_options = "SOURCE_DELETION_NEVER"
    authentication {
      sqs_access_key_id = "XXXXX"	
      sqs_secret_access_key = "XXXXX"

      s3_authentication {
        access_key_id = "XXXXX"	
        secret_access_key = "XXXXX"
      }
    }
  }

}
