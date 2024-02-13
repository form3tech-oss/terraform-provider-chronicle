package client

const (
	sqsFeedConfigurationPropertyKey = "amazonSqsSettings"
)

type SQSFeedConfiguration struct {
	Queue               string                `json:"queue,omitempty"`
	Region              string                `json:"region,omitempty"`
	AccountNumber       string                `json:"accountNumber,omitempty"`
	SourceDeleteOptions string                `json:"sourceDeletionOption,omitempty"`
	Authentication      SQSFeedAuthentication `json:"authentication,omitempty"`
}
type SQSFeedAuthentication struct {
	SQSAuthentication SQSFeedAuthenticationCred  `json:"sqsAccessKeySecretAuth,omitempty"`
	S3Authentication  *SQSFeedAuthenticationCred `json:"additionalS3AccessKeySecretAuth,omitempty"`
}
type SQSFeedAuthenticationCred struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

func (s3 *SQSFeedConfiguration) getConfigurationPropertyKey() string {
	return sqsFeedConfigurationPropertyKey
}
func (s3 *SQSFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeSQS
}
