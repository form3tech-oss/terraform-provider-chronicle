package client

const (
	AzureBlobStoreV2FeedConfigurationPropertyKey = "azureBlobStoreV2Settings"
)

type AzureBlobStoreV2FeedConfiguration struct {
	AzureURI            string                                 `json:"azureUri,omitempty"`
	SourceDeleteOptions string                                 `json:"sourceDeletionOption,omitempty"`
	MaxLookbackDays     int                                    `json:"maxLookbackDays,omitempty"`
	Authentication      AzureBlobStoreV2FeedFeedAuthentication `json:"authentication,omitempty"`
}

type AzureBlobStoreV2FeedFeedAuthentication struct {
	SharedKey string `json:"sharedKey,omitempty"`
	SASToken  string `json:"sasToken,omitempty"`
}

func (a *AzureBlobStoreV2FeedConfiguration) getConfigurationPropertyKey() string {
	return AzureBlobStoreV2FeedConfigurationPropertyKey
}

func (a *AzureBlobStoreV2FeedConfiguration) getFeedSourceType() string {
	return FeedSourceTypeAzureBlobStoreV2
}
