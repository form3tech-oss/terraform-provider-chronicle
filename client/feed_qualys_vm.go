package client

const (
	qualysVMFeedConfigurationPropertyKey = "qualysVmSettings"
	QualysVMFeedLogType                  = "QUALYS_VM"
)

type QualysVMFeedConfiguration struct {
	Hostname       string                     `json:"hostname,omitempty"`
	Authentication QualysVMFeedAuthentication `json:"authentication,omitempty"`
}
type QualysVMFeedAuthentication struct {
	User   string `json:"user,omitempty"`
	Secret string `json:"secret,omitempty"`
}

func (s3 *QualysVMFeedConfiguration) getConfigurationPropertyKey() string {
	return qualysVMFeedConfigurationPropertyKey
}
func (s3 *QualysVMFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAPI
}
