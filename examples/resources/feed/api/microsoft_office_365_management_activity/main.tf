resource "chronicle_feed_microsoft_office_365_management_activity" "feed" {
  display_name = "officefeed"
  enabled      = false
  namespace    = "one"
  labels = {
    "env" = "one"
  }
  details {
    hostname     = "manage.office.com/api/v1.0"
    tenant_id    = "XXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXX"
    content_type = "AUDIT_AZURE_ACTIVE_DIRECTORY"
    authentication {
      client_id     = "XXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXX"
      client_secret = "XXXXXXXX"
    }
  }

}
