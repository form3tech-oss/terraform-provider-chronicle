resource "chronicle_feed_qualys_vm" "feed" {
  display_name = "qualysvmfeed"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    hostname = "qualysapi.qualys.com/api/2.0/fo/asset/host/?action=list"
    authentication {
      user   = "XXXXX"
      secret = "XXXX"
    }
  }

}
