package client

const (
	s3V2FeedConfigurationPropertyKey = "amazonS3V2Settings"
)

type S3V2FeedConfiguration struct {
	S3URI               string                 `json:"s3Uri,omitempty"`
	SourceDeleteOptions string                 `json:"sourceDeletionOption,omitempty"`
	MaxLookbackDays     int                    `json:"maxLookbackDays,omitempty"`
	Authentication      S3V2FeedAuthentication `json:"authentication,omitempty"`
}

type S3V2FeedAuthentication struct {
	Region          string `json:"region,omitempty"`
	AccessKeyID     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

func (s *S3V2FeedConfiguration) getConfigurationPropertyKey() string {
	return s3V2FeedConfigurationPropertyKey
}

func (s *S3V2FeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeS3V2
}
