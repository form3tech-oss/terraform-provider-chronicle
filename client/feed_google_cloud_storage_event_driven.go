package client

const (
	GCSEventDrivenFeedConfigurationPropertyKey = "googleCloudStorageEventDrivenSettings"
)

type GCSEventDrivenFeedConfiguration struct {
	BucketURI           string `json:"bucketUri,omitempty"`
	PubsubSubscription  string `json:"pubsubSubscription,omitempty"`
	SourceDeleteOptions string `json:"sourceDeletionOption,omitempty"`
	MaxLookbackDays     int    `json:"maxLookbackDays,omitempty"`
}

func (g *GCSEventDrivenFeedConfiguration) getConfigurationPropertyKey() string {
	return GCSEventDrivenFeedConfigurationPropertyKey
}

func (g *GCSEventDrivenFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeGCSEventDriven
}
