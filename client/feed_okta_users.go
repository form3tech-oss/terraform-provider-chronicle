package client

const (
	oktaUsersFeedConfigurationPropertyKey = "oktaUserContextSettings"
	OktaUsersFeedLogType                  = "OKTA_USER_CONTEXT"
)

type OktaUsersFeedConfiguration struct {
	Hostname                string                      `json:"hostname,omitempty"`
	ManagerIDReferenceField string                      `json:"managerIdReferenceField,omitempty"`
	Authentication          OktaUsersFeedAuthentication `json:"authentication,omitempty"`
}
type OktaUsersFeedAuthentication struct {
	HeaderKeyValues []OktaUsersAuthenticationHeaderKeyValues `json:"headerKeyValues,omitempty"`
}

type OktaUsersAuthenticationHeaderKeyValues struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (s3 *OktaUsersFeedConfiguration) getConfigurationPropertyKey() string {
	return oktaUsersFeedConfigurationPropertyKey
}
func (s3 *OktaUsersFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAPI
}
