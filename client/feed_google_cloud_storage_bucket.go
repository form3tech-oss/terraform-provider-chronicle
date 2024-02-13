package client

const (
	GCPBucketFeedConfigurationPropertyKey = "gcsSettings"
)

type GCPBucketFeedConfiguration struct {
	URI                 string `json:"bucketUri,omitempty"`
	SourceType          string `json:"sourceType,omitempty"`
	SourceDeleteOptions string `json:"sourceDeletionOption,omitempty"`
}

func (b *GCPBucketFeedConfiguration) getConfigurationPropertyKey() string {
	return GCPBucketFeedConfigurationPropertyKey
}
func (b *GCPBucketFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeGCS
}
