---
page_title: "Provider: Chronicle"
description: |-
  The Chronicle provider.
---
# chronicle Provider

## Example Usage

```terraform
provider "chronicle" {
  backstoryapi_credentials = "backstoryapi_crendentials_json_string"
  region                   = "europe"
  request_attempts         = 5
  request_timeout          = 120
}
```


## Configuration
Note that for each API you can only provide either credentials or access token. Environment variables always take the lowest precedence.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `alert_custom_endpoint` (String) Custom URL to alert endpoint.
- `alias_custom_endpoint` (String) Custom URL to alias endpoint.
- `artifact_custom_endpoint` (String) Custom URL to artifact endpoint.
- `asset_custom_endpoint` (String) Custom URL to asset endpoint.
- `backstoryapi_access_token` (String) Backstory API access token. Local file path or content.
- `backstoryapi_credentials` (String) Backstory API credential. Local file path or content.
				 It may be replaced by CHRONICLE_BACKSTORY_CREDENTIALS environment variable, which expects base64 encoded credential.
- `bigqueryapi_access_token` (String) BigQuery API access token. Local file path or content.
- `bigqueryapi_credentials` (String) BigQuery API crendential. Local file path or content.
				 It may be replaced by CHRONICLE_BIGQUERY_CREDENTIALS environment variable, which expects base64 encoded credential.
- `events_custom_endpoint` (String) Custom URL to events endpoint.
- `feed_custom_endpoint` (String) Custom URL to feed endpoint.
- `forwarderapi_access_token` (String) Forwarder API Access token. Local file path or content.
- `forwarderapi_credentials` (String) Forwarder API crendential. Local file path or content.
				 It may be replaced by CHRONICLE_FORWARDER_CREDENTIALS environment variable, which expects base64 encoded credential.
- `ingestionapi_access_token` (String) Ingestion API access token. Local file path or content.
- `ingestionapi_credentials` (String) Ingestion API crendential. Local file path or content.
				 It may be replaced by CHRONICLE_INGESTION_CREDENTIALS environment variable, which expects base64 encoded credential.
- `ioc_custom_endpoint` (String) Custom URL to ioc endpoint.
- `region` (String) Region to which send requests, available regions are: [us europe europe-west2 asia-southeast1]. It may be replaced by CHRONICLE_REGION environment variable.
- `request_attempts` (Number) Number of attempts per request. Attempts follow exponential back-off strategy. Defaults to 5 attempts.
- `request_timeout` (Number) Request timeout in seconds. Defaults to 120 (s).
- `rule_custom_endpoint` (String) Custom URL to rule endpoint.
- `subjects_custom_endpoint` (String) Custom URL to subjects endpoint.
