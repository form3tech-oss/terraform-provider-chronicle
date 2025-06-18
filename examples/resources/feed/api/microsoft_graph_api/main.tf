// Defaults to provider subscription id
data "azurerm_subscription" "current" {
}
// Use current principal to own azuread_application resource
data "azurerm_client_config" "current" {}

locals {
  resource_tags = {
    "managed-by" = "terraform"
  }
  event_namespace = "my-namespace"
}

// Create Application servicePrincipal with rights to Graph API
resource "azuread_application" "chronicle_graphapi" {
  display_name = "chronicle-graphapi-ingest"
  owners       = [data.azuread_client_config.current.object_id]
  required_resource_access {
    # Retrieved with:
    # az ad sp list --display-name "Microsoft Graph" --query '[].{appDisplayName:appDisplayName, appId:appId}'
    resource_app_id = "00000003-0000-0000-c000-000000000000" # Microsoft Graph
    # Retrieved with:
    # az ad sp show --id 00000003-0000-0000-c000-000000000000 | \
    #   jq '.appRoles[] | select (.value == "SecurityActions.Read.All" or .value == "SecurityEvens.Read.All")'
    resource_access {
      id   = "5e0edab9-c148-49d0-b423-ac253e121825" # SecurityActions.Read.All
      type = "Role"
    }
    resource_access {
      id   = "bf394140-e372-4bf9-a898-299cfc7564e5" # SecurityEvents.Read.All
      type = "Role"
    }
    resource_access {
      id   = "b0afded3-3588-46d8-8b3d-9842eff778da" # AuditLog.Read.All
      type = "Role"
    }
    resource_access {
      id   = "7ab1d382-f21e-4acd-a863-ba3e13f7da61" # Directory.Read.All
      type = "Role"
    }
  }
  tags = [for k, v in local.resource_tags : "${k}:${v}"]
}

// Create an application password
resource "azuread_application_password" "graphapi_pw" {
  display_name   = "Chronicle GraphAPI Feeds"
  application_id = azuread_application.chronicle_graphapi.id
}

// Graph Alerts
resource "chronicle_feed_microsoft_graph_api" "alerts" {
  enabled      = true
  display_name = "MS Graph - Alerts"
  namespace    = local.event_namespace
  log_type     = "MICROSOFT_GRAPH_ALERT"
  details {
    tenant_id     = data.azurerm_subscription.current.tenant_id
    hostname      = "graph.microsoft.com/v1.0/security/alerts"
    auth_endpoint = "login.microsoftonline.com"
    authentication {
      client_id     = azuread_application.chronicle_graphapi.client_id
      client_secret = azuread_application_password.graphapi_pw.value
    }
  }
  labels = merge(local.resource_tags, { "ingest-path" = "api"})
}

// Capture directoryAudit events
resource "chronicle_feed_microsoft_graph_api" "azure_ad_audit" {
  enabled      = true
  display_name = "MS Graph - directoryAudit"
  namespace    = local.event_namespace
  log_type     = "AZURE_AD_AUDIT"
  details {
    tenant_id = data.azurerm_subscription.current.tenant_id
    hostname  = "graph.microsoft.com/v1.0/auditLogs/directoryAudits"
    authentication {
      client_id     = azuread_application.chronicle_graphapi.client_id
      client_secret = azuread_application_password.graphapi_pw.value
    }
  }
  labels = merge(local.resource_tags, { "ingest-path" = "api" })
}

// Sign-In events
resource "chronicle_feed_microsoft_graph_api" "azure_ad" {
  enabled      = true
  display_name = "MS Graph - signIns"
  namespace    = local.event_namespace
  log_type     = "AZURE_AD"
  details {
    tenant_id = data.azurerm_subscription.current.tenant_id
    hostname  = "graph.microsoft.com/v1.0/auditLogs/signIns"
    authentication {
      client_id     = azuread_application.chronicle_graphapi.client_id
      client_secret = azuread_application_password.graphapi_pw.value
    }
  }
  labels = merge(local.resource_tags, { "ingest-path" = "api" })
}

// Use Azure AD for context enrichment
resource "chronicle_feed_microsoft_graph_api" "azure_ad_context" {
  enabled      = true
  display_name = "MS Graph - AD Context"
  namespace    = local.event_namespace
  log_type     = "AZURE_AD_CONTEXT"
  details {
    tenant_id        = data.azurerm_subscription.current.tenant_id
    hostname         = "graph.microsoft.com/beta"
    retrieve_devices = true
    retrieve_groups  = true
    authentication {
      client_id     = azuread_application.chronicle_graphapi.client_id
      client_secret = azuread_application_password.graphapi_pw.value
    }
  }
  labels = merge(local.resource_tags, { "ingest-path" = "api" })
}
