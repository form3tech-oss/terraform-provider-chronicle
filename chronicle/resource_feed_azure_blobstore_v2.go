package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedAzureBlobStoreV2 struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedAzureBlobStoreV2() *ResourceFeedAzureBlobStoreV2 {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"azure_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The Azure Blob Storage URI in the format https://account.blob.core.windows.net/container/path/.`,
			},
			"source_delete_options": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateFeedV2SourceDeleteOption,
				Required:         true,
				Description: `Whether to delete source files after they have been transferred to Chronicle. Valid values are:

- NEVER: Never delete files from the source.
- ON_SUCCESS: Delete files and empty directories from the source after successful ingestion.`,
			},
			"max_lookback_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     180,
				Description: `The maximum number of days in the past to look for files. Default is 180 days.`,
			},
			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Azure authentication details. Specify either shared_key or sas_token, but not both.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shared_key": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							ConflictsWith: []string{"details.0.authentication.0.sas_token"},
							Description:   `The shared access key for the Azure Blob Storage account.`,
						},
						"sas_token": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							ConflictsWith: []string{"details.0.authentication.0.shared_key"},
							Description:   `The SAS (Shared Access Signature) token for the Azure Blob Storage account.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a V2 feed from Azure Blob Storage. This feed type uses the Google Cloud Storage Transfer Service for improved ingestion."
	resource := &ResourceFeedAzureBlobStoreV2{}
	resource.TerraformResource = newFeedResourceSchema(details, resource, description, true)

	return resource
}

func (f *ResourceFeedAzureBlobStoreV2) getLogType() string {
	return ""
}

func (f *ResourceFeedAzureBlobStoreV2) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	configuration := &chronicle.AzureBlobStoreV2FeedConfiguration{
		AzureURI:            resourceDetails["azure_uri"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		MaxLookbackDays:     resourceDetails["max_lookback_days"].(int),
		Authentication:      chronicle.AzureBlobStoreV2FeedFeedAuthentication{},
	}

	if sharedKey, ok := authenticationDetails["shared_key"].(string); ok && sharedKey != "" {
		configuration.Authentication.SharedKey = sharedKey
	}

	if sasToken, ok := authenticationDetails["sas_token"].(string); ok && sasToken != "" {
		configuration.Authentication.SASToken = sasToken
	}

	return configuration
}

//nolint:all
func (f *ResourceFeedAzureBlobStoreV2) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readAzureConf := readConf.(*chronicle.AzureBlobStoreV2FeedConfiguration)

	// Import Case
	if originalConf == nil {
		authMap := make(map[string]interface{})
		if readAzureConf.Authentication.SharedKey != "" {
			authMap["shared_key"] = readAzureConf.Authentication.SharedKey
		}
		if readAzureConf.Authentication.SASToken != "" {
			authMap["sas_token"] = readAzureConf.Authentication.SASToken
		}

		return []map[string]interface{}{{
			"azure_uri":             readAzureConf.AzureURI,
			"source_delete_options": readAzureConf.SourceDeleteOptions,
			"max_lookback_days":     readAzureConf.MaxLookbackDays,
			"authentication":        []map[string]interface{}{authMap},
		}}
	}

	originalAzureConf := originalConf.(*chronicle.AzureBlobStoreV2FeedConfiguration)
	// Default Case
	authMap := make(map[string]interface{})
	if originalAzureConf.Authentication.SharedKey != "" {
		authMap["shared_key"] = originalAzureConf.Authentication.SharedKey
	}
	if originalAzureConf.Authentication.SASToken != "" {
		authMap["sas_token"] = originalAzureConf.Authentication.SASToken
	}

	return []map[string]interface{}{{
		"azure_uri":             readAzureConf.AzureURI,
		"source_delete_options": originalAzureConf.SourceDeleteOptions, // not returned
		"max_lookback_days":     readAzureConf.MaxLookbackDays,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{authMap},
	}}
}
