resource "chronicle_feed_proofpoint_siem" "test" {
  display_name = "ppsiem"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    authentication {
      user   = "XXXXX"
      secret = "XXXX"
    }
  }
}
