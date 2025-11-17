package client

const (
	GCSV2FeedConfigurationPropertyKey = "gcsV2Settings"
)

type GCSV2FeedConfiguration struct {
	BucketURI           string `json:"bucketUri,omitempty"`
	SourceDeleteOptions string `json:"sourceDeletionOption,omitempty"`
	MaxLookbackDays     int    `json:"maxLookbackDays,omitempty"`
}

func (g *GCSV2FeedConfiguration) getConfigurationPropertyKey() string {
	return GCSV2FeedConfigurationPropertyKey
}

func (g *GCSV2FeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeGCSV2
}
