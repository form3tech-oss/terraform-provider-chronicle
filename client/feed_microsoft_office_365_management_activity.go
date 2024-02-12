package client

const (
	microsoftOffice365ManagementActivityFeedConfigurationPropertyKey = "office365Settings"
	MicrosoftOffice365ManagementActivityFeedLogType                  = "OFFICE_365"
)

type MicrosoftOffice365ManagementActivityFeedConfiguration struct {
	TenantID       string                                                 `json:"tenantId,omitempty"`
	ContentType    string                                                 `json:"contentType,omitempty"`
	Hostname       string                                                 `json:"hostname,omitempty"`
	Authentication MicrosoftOffice365ManagementActivityFeedAuthentication `json:"authentication,omitempty"`
}
type MicrosoftOffice365ManagementActivityFeedAuthentication struct {
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

func (s3 *MicrosoftOffice365ManagementActivityFeedConfiguration) getConfigurationPropertyKey() string {
	return microsoftOffice365ManagementActivityFeedConfigurationPropertyKey
}
func (s3 *MicrosoftOffice365ManagementActivityFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAPI
}
