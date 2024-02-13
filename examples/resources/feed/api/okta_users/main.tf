resource "chronicle_feed_okta_users" "test" {
  display_name = "example"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    hostname   = "something.okta.com"
    manager_id = "123"
    authentication {
      key   = "Authentication"
      value = "123"
    }
  }
}
