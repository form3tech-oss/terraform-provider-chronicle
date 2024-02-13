package client

const (
	s3FeedConfigurationPropertyKey = "amazonS3Settings"
)

type S3FeedConfiguration struct {
	URI                 string               `json:"s3Uri,omitempty"`
	SourceType          string               `json:"sourceType,omitempty"`
	SourceDeleteOptions string               `json:"sourceDeletionOption,omitempty"`
	Authentication      S3FeedAuthentication `json:"authentication,omitempty"`
}
type S3FeedAuthentication struct {
	Region          string `json:"region,omitempty"`
	AccessKeyID     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

func (s3 *S3FeedConfiguration) getConfigurationPropertyKey() string {
	return s3FeedConfigurationPropertyKey
}
func (s3 *S3FeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeS3
}
