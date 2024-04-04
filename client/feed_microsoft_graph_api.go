package client

type MicrosoftGraphAPIFeedConfiguration struct {
	TenantID       string                              `json:"tenantId,omitempty"`
	ContentType    string                              `json:"-"`
	Hostname       string                              `json:"hostname,omitempty"`
	Authentication MicrosoftGraphAPIFeedAuthentication `json:"authentication,omitempty"`

	// Only for AZURE_AD_CONTEXT
	RetrieveDevices bool `json:"retrieveDevices,omitempty"`
	RetrieveGroups  bool `json:"retrieveGroups,omitempty"`

	// Only for MICROSOFT_GRAPH_ALERT
	AuthEndpoint string `json:"authEndpoint,omitempty"`
}
type MicrosoftGraphAPIFeedAuthentication struct {
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

func (config *MicrosoftGraphAPIFeedConfiguration) getConfigurationPropertyKey() string {
	switch key := config.ContentType; key {
	case "AZURE_AD_AUDIT":
		return "azureAdAuditSettings"
	case "AZURE_AD_CONTEXT":
		return "azureAdContextSettings"
	case "AZURE_AD":
		return "azureAdSettings"
	case "AZURE_MDM_INTUNE":
		return "azureMdmIntuneSettings"
	case "MICROSOFT_GRAPH_ALERT":
		return "microsoftGraphAlertSettings"
	default:
		return ""
	}
}
func (config *MicrosoftGraphAPIFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAPI
}
