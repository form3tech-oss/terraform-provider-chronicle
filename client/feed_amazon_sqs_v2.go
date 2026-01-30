package client

const (
	sqsV2FeedConfigurationPropertyKey = "amazon_sqs_v2_settings"
)

type SQSV2FeedConfiguration struct {
	Queue               string                  `json:"queue,omitempty"`
	S3URI               string                  `json:"s3Uri,omitempty"`
	SourceDeleteOptions string                  `json:"sourceDeletionOption,omitempty"`
	MaxLookbackDays     int                     `json:"maxLookbackDays,omitempty"`
	Authentication      SQSV2FeedAuthentication `json:"authentication,omitempty"`
}

type SQSV2FeedAuthentication struct {
	AccessKeySecretAuth *SQSV2AccessKeySecretAuth `json:"access_key_secret_auth,omitempty"`
	AWSIAMRoleArn       string                    `json:"aws_iam_role_arn,omitempty"`
}

type SQSV2AccessKeySecretAuth struct {
	AccessKeyID     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

func (s *SQSV2FeedConfiguration) getConfigurationPropertyKey() string {
	return sqsV2FeedConfigurationPropertyKey
}

func (s *SQSV2FeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeSQSV2
}
