resource "chronicle_rule" "test" {
  rule_text        = file("path/to/yararule")
  alerting_enabled = false
  live_enabled     = false
}
