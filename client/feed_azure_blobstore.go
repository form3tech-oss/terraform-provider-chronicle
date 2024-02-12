package client

const (
	AzureBlobStoreFeedConfigurationPropertyKey = "azureBlobStoreSettings"
)

type AzureBlobStoreFeedConfiguration struct {
	URI                 string                               `json:"azureUri,omitempty"`
	SourceType          string                               `json:"sourceType,omitempty"`
	SourceDeleteOptions string                               `json:"sourceDeletionOption,omitempty"`
	Authentication      AzureBlobStoreFeedFeedAuthentication `json:"authentication,omitempty"`
}
type AzureBlobStoreFeedFeedAuthentication struct {
	SharedKey string `json:"sharedKey,omitempty"`
	SASToken  string `json:"sasToken,omitempty"`
}

func (b *AzureBlobStoreFeedConfiguration) getConfigurationPropertyKey() string {
	return AzureBlobStoreFeedConfigurationPropertyKey
}
func (b *AzureBlobStoreFeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAzureBlobStore
}
