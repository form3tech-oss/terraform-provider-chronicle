package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionNever              = "SOURCE_DELETION_NEVER"
	FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionOnSuccess          = "SOURCE_DELETION_ON_SUCCESS"
	FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionOnSuccessFilesOnly = "SOURCE_DELETION_ON_SUCCESS_FILES_ONLY"
)

const (
	FeedGoogleCloudStorageBucketSourceTypeFiles            = "FILES"
	FeedGoogleCloudStorageBucketSourceTypeFolders          = "FOLDERS"
	FeedGoogleCloudStorageBucketSourceTypeFoldersRecursive = "FOLDERS_RECURSIVE"
)

type ResourceFeedGoogleCloudStorageBucket struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedGoogleCloudStorageBucket() *ResourceFeedGoogleCloudStorageBucket {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_uri": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateGCSURI,
				Required:         true,
				Description:      `The bucket URI to ingest.`,
			},
			"bucket_source_type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateFeedGCSSourceType,
				Required:         true,
				Description: `The type of file indicated by the uri. It may be the following:

				- FILES: The URI points to a single file which will be ingested with each execution of the feed.
				- FOLDERS: The URI points to a directory. All files contained within the directory will be ingested with each execution of the feed.
				- FOLDERS_RECURSIVE: The URI points to a directory. All files and directories contains within the indicated directory will be ingested,
				 including all files and directories within those directories, and so on. `,
			},

			"source_delete_options": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateFeedGCSSourceDeleteOption,
				Required:         true,
				Description: `Whether to delete source files after they have been transferred to Chronicle. The possible values are as follows:

				- SOURCE_DELETION_NEVER: Never delete files from the source.
				- SOURCE_DELETION_ON_SUCCESS: Delete files and empty directories from the source after successful ingestion.
				- SOURCE_DELETION_ON_SUCCESS_FILES_ONLY: Delete files from the source after successful ingestion.`,
			},
		},
	}
	description := "Creates a feed from Google Cloud Storage Bucket service."
	bucket := &ResourceFeedGoogleCloudStorageBucket{}
	bucket.TerraformResource = newFeedResourceSchema(details, bucket, description, true)

	return bucket
}

func (f *ResourceFeedGoogleCloudStorageBucket) getLogType() string {
	return ""
}

func (f *ResourceFeedGoogleCloudStorageBucket) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})

	return &chronicle.GCPBucketFeedConfiguration{
		URI:                 resourceDetails["bucket_uri"].(string),
		SourceType:          resourceDetails["bucket_source_type"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
	}
}

//nolint:all
func (f *ResourceFeedGoogleCloudStorageBucket) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {
	readBucketConf := readConf.(*chronicle.GCPBucketFeedConfiguration)

	return []map[string]interface{}{{
		"bucket_uri":            readBucketConf.URI,
		"bucket_source_type":    readBucketConf.SourceType,
		"source_delete_options": readBucketConf.SourceDeleteOptions,
	}}
}
