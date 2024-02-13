resource "chronicle_feed_thinkst_canary" "test" {
  display_name = "example"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    hostname = "something.canary.tools"
    authentication {
      key   = "auth_token"
      value = "123"
    }
  }
}
