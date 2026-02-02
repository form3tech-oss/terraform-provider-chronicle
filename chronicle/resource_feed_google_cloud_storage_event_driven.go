package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedGoogleCloudStorageEventDriven struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedGoogleCloudStorageEventDriven() *ResourceFeedGoogleCloudStorageEventDriven {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_uri": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateGCSURI,
				Description:      `The Google Cloud Storage bucket URI in the format gs://bucket-name/path/.`,
			},
			"pubsub_subscription": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The Pub/Sub subscription name. Format: projects/your-project/subscriptions/your-subscription`,
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
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          180,
				ValidateDiagFunc: validateMaxLookbackDays,
				Description:      `The maximum number of days in the past to look for files. Must be between 1 and 180. Default is 180 days.`,
			},
		},
	}
	description := "Creates an event-driven feed from Google Cloud Storage using Pub/Sub notifications. " +
		"This feed type uses Google Cloud Storage Transfer Service with push-based ingestion for reduced latency. " +
		"Authentication is handled via the Google Security Operations service account."
	resource := &ResourceFeedGoogleCloudStorageEventDriven{}
	resource.TerraformResource = newFeedResourceSchema(details, resource, description, true)

	return resource
}

func (f *ResourceFeedGoogleCloudStorageEventDriven) getLogType() string {
	return ""
}

func (f *ResourceFeedGoogleCloudStorageEventDriven) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})

	return &chronicle.GCSEventDrivenFeedConfiguration{
		BucketURI:           resourceDetails["bucket_uri"].(string),
		PubsubSubscription:  resourceDetails["pubsub_subscription"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		MaxLookbackDays:     resourceDetails["max_lookback_days"].(int),
	}
}

//nolint:all
func (f *ResourceFeedGoogleCloudStorageEventDriven) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readGCSConf := readConf.(*chronicle.GCSEventDrivenFeedConfiguration)

	// Import Case
	if originalConf == nil {
		return []map[string]interface{}{{
			"bucket_uri":            readGCSConf.BucketURI,
			"pubsub_subscription":   readGCSConf.PubsubSubscription,
			"source_delete_options": readGCSConf.SourceDeleteOptions,
			"max_lookback_days":     readGCSConf.MaxLookbackDays,
		}}
	}

	originalGCSConf := originalConf.(*chronicle.GCSEventDrivenFeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		"bucket_uri":            readGCSConf.BucketURI,
		"pubsub_subscription":   readGCSConf.PubsubSubscription,
		"source_delete_options": originalGCSConf.SourceDeleteOptions, // not returned
		"max_lookback_days":     readGCSConf.MaxLookbackDays,
	}}
}
