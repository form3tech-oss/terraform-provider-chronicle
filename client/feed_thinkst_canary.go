package client

const (
	thinkstCanaryFeedConfigurationPropertyKey = "thinkstCanarySettings"
	ThinkstCanaryFeedLogType                  = "THINKST_CANARY"
)

type ThinkstCanaryFeedConfiguration struct {
	Hostname       string                          `json:"hostname,omitempty"`
	Authentication ThinkstCanaryFeedAuthentication `json:"authentication,omitempty"`
}

type ThinkstCanaryFeedAuthentication struct {
	HeaderKeyValues []ThinkstCanaryAuthenticationHeaderKeyValues `json:"headerKeyValues,omitempty"`
}

type ThinkstCanaryAuthenticationHeaderKeyValues struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (c *ThinkstCanaryFeedConfiguration) getConfigurationPropertyKey() string {
	return thinkstCanaryFeedConfigurationPropertyKey
}
func (c *ThinkstCanaryFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAPI
}
