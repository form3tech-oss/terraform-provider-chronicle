package client

const (
	oktaSystemLogFeedConfigurationPropertyKey = "oktaSettings"
	OktaSystemLogFeedLogType                  = "OKTA"
)

type OktaSystemLogFeedConfiguration struct {
	Hostname       string                          `json:"hostname,omitempty"`
	Authentication OktaSystemLogFeedAuthentication `json:"authentication,omitempty"`
}
type OktaSystemLogFeedAuthentication struct {
	HeaderKeyValues []OktaSystemLogAuthenticationHeaderKeyValues `json:"headerKeyValues,omitempty"`
}

type OktaSystemLogAuthenticationHeaderKeyValues struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (c *OktaSystemLogFeedConfiguration) getConfigurationPropertyKey() string {
	return oktaSystemLogFeedConfigurationPropertyKey
}
func (c *OktaSystemLogFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAPI
}
