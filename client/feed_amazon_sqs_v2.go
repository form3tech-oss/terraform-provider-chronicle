package client

const (
	sqsV2FeedConfigurationPropertyKey = "amazonSqsV2Settings"
)

type SQSV2FeedConfiguration struct {
	S3URI               string                  `json:"s3Uri,omitempty"`
	Region              string                  `json:"region,omitempty"`
	AccountNumber       string                  `json:"accountNumber,omitempty"`
	QueueName           string                  `json:"queueName,omitempty"`
	SourceDeleteOptions string                  `json:"sourceDeletionOption,omitempty"`
	MaxLookbackDays     int                     `json:"maxLookbackDays,omitempty"`
	Authentication      SQSV2FeedAuthentication `json:"authentication,omitempty"`
}

type SQSV2FeedAuthentication struct {
	SQSAuthentication SQSV2FeedAuthenticationCred  `json:"sqsAccessKeySecretAuth,omitempty"`
	S3Authentication  *SQSV2FeedAuthenticationCred `json:"additionalS3AccessKeySecretAuth,omitempty"`
}

type SQSV2FeedAuthenticationCred struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

func (s *SQSV2FeedConfiguration) getConfigurationPropertyKey() string {
	return sqsV2FeedConfigurationPropertyKey
}

func (s *SQSV2FeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeSQSV2
}
