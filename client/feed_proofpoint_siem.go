package client

const (
	proofpointSIEMFeedConfigurationPropertyKey = "proofpointMailSettings"
	ProofpointSIEMFeedLogType                  = "PROOFPOINT_MAIL"
)

type ProofpointSIEMFeedConfiguration struct {
	Authentication ProofpointSIEMFeedAuthentication `json:"authentication,omitempty"`
}
type ProofpointSIEMFeedAuthentication struct {
	User   string `json:"user,omitempty"`
	Secret string `json:"secret,omitempty"`
}

func (s3 *ProofpointSIEMFeedConfiguration) getConfigurationPropertyKey() string {
	return proofpointSIEMFeedConfigurationPropertyKey
}
func (s3 *ProofpointSIEMFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAPI
}
