package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedGoogleCloudStorageV2 struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedGoogleCloudStorageV2() *ResourceFeedGoogleCloudStorageV2 {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_uri": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateGCSURI,
				Description:      `The Google Cloud Storage bucket URI in the format gs://bucket-name/path/.`,
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
		},
	}
	description := "Creates a V2 feed from Google Cloud Storage. This feed type uses the Google Cloud Storage Transfer Service for improved ingestion. Authentication is handled via the Google Security Operations service account."
	resource := &ResourceFeedGoogleCloudStorageV2{}
	resource.TerraformResource = newFeedResourceSchema(details, resource, description, true)

	return resource
}

func (f *ResourceFeedGoogleCloudStorageV2) getLogType() string {
	return ""
}

func (f *ResourceFeedGoogleCloudStorageV2) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})

	return &chronicle.GCSV2FeedConfiguration{
		BucketURI:           resourceDetails["bucket_uri"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		MaxLookbackDays:     resourceDetails["max_lookback_days"].(int),
	}
}

//nolint:all
func (f *ResourceFeedGoogleCloudStorageV2) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readGCSConf := readConf.(*chronicle.GCSV2FeedConfiguration)

	// Import Case
	if originalConf == nil {
		return []map[string]interface{}{{
			"bucket_uri":            readGCSConf.BucketURI,
			"source_delete_options": readGCSConf.SourceDeleteOptions,
			"max_lookback_days":     readGCSConf.MaxLookbackDays,
		}}
	}

	originalGCSConf := originalConf.(*chronicle.GCSV2FeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		"bucket_uri":            readGCSConf.BucketURI,
		"source_delete_options": originalGCSConf.SourceDeleteOptions, // not returned
		"max_lookback_days":     readGCSConf.MaxLookbackDays,
	}}
}
