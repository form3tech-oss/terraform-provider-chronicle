resource "chronicle_azure_blobstore" "azure" {
  display_name = "mys3feed"
  log_type     = "GITHUB"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    uri = "https://myaccount.blob.core.windows.net/logging"
    source_type = "FOLDERS_RECURSIVE"
    authentication {
      shared_key = "XXXX"
    }
  }

}
