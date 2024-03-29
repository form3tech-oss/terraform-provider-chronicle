---
page_title: "chronicle_feed_microsoft_office_365_management_activity Resource - terraform-provider-chronicle"
subcategory: ""
description: |-
  Creates a feed from API source type for Microsoft Office 365 Management Activity log type.
---

# chronicle_feed_microsoft_office_365_management_activity (Resource)

Creates a feed from API source type for Microsoft Office 365 Management Activity log type.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `details` (Block List, Min: 1, Max: 1) Each feed type has its own requirements for which this field must fulfil. (see [below for nested schema](#nestedblock--details))
- `display_name` (String) Name to be displayed.
- `enabled` (Boolean) Enabled specifies whether a feed is allowed to be executed.

### Optional

- `labels` (Map of String) All of the events that result from this feed will have this label applied.
- `namespace` (String) The namespace the feed will be associated with.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `feed_source_type` (String) Feed Source Type describes how data is collected.
- `id` (String) The ID of this resource.
- `log_type` (String) Log Type is a label which describes the nature of the data being ingested.
- `state` (String) State gives some insight into the current state of a feed.

<a id="nestedblock--details"></a>
### Nested Schema for `details`

Required:

- `authentication` (Block List, Min: 1, Max: 1) Office 365 authentication details. (see [below for nested schema](#nestedblock--details--authentication))
- `content_type` (String) The type of logs to fetch. See https://cloud.google.com/chronicle/docs/reference/feed-management-api#office_365_content_type.
- `tenant_id` (String) Tenant ID (a UUID).

Optional:

- `hostname` (String) API Full Path, default value: manage.office.com/api/v1.0.

<a id="nestedblock--details--authentication"></a>
### Nested Schema for `details.authentication`

Required:

- `client_id` (String) OAuth client ID (a UUID).
- `client_secret` (String, Sensitive) OAuth client secret.



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)
