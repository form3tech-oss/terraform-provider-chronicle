resource "chronicle_feed_okta_system_log" "test" {
  display_name = "example"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    hostname = "something.okta.com"
    authentication {
      key   = "Authentication"
      value = "123"
    }
  }
}
