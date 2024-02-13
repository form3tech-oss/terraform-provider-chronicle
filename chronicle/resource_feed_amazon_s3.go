package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FeedS3SourceDeleteOptionDeletionNever              = "SOURCE_DELETION_NEVER"
	FeedS3SourceDeleteOptionDeletionOnSuccess          = "SOURCE_DELETION_ON_SUCCESS"
	FeedS3SourceDeleteOptionDeletionOnSuccessFilesOnly = "SOURCE_DELETION_ON_SUCCESS_FILES_ONLY"
)

const (
	FeedS3SourceTypeFiles            = "FILES"
	FeedS3SourceTypeFolders          = "FOLDERS"
	FeedS3SourceTypeFoldersRecursive = "FOLDERS_RECURSIVE"
)

type ResourceFeedAmazonS3 struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedAmazonS3() *ResourceFeedAmazonS3 {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"s3_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The S3 URI to ingest.`,
			},
			"s3_source_type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateFeedS3SourceType,
				Required:         true,
				Description: `The type of file indicated by the uri. It may be the following:

				- FILES: The URI points to a single file which will be ingested with each execution of the feed.
				- FOLDERS: The URI points to a directory. All files contained within the directory will be ingested with each execution of the feed.
				- FOLDERS_RECURSIVE: The URI points to a directory. All files and directories contains within the indicated directory will be ingested,
				 including all files and directories within those directories, and so on. `,
			},

			"source_delete_options": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateFeedS3SourceDeleteOption,
				Required:         true,
				Description: `Whether to delete source files after they have been transferred to Chronicle. The possible values are as follows:

				- SOURCE_DELETION_NEVER: Never delete files from the source.
				- SOURCE_DELETION_ON_SUCCESS: Delete files and empty directories from the source after successful ingestion.
				- SOURCE_DELETION_ON_SUCCESS_FILES_ONLY: Delete files from the source after successful ingestion.`,
			},

			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `AWS authentication details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The region where the S3 bucket resides following this format: https://cloud.google.com/chronicle/docs/reference/feed-management-api#amazon_s3_regions.`,
						},
						"access_key_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateAWSAccessKeyID,
							Description:      `This is the 20 character ID associated with your Amazon IAM account.`,
						},
						"secret_access_key": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validateAWSSecretAccessKey,
							Description:      `This is the 40 character access key associated with your Amazon IAM account.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a feed from Amazon Simple Storage Service Bucket."
	s3 := &ResourceFeedAmazonS3{}
	s3.TerraformResource = newFeedResourceSchema(details, s3, description, true)

	return s3
}

func (f *ResourceFeedAmazonS3) getLogType() string {
	return ""
}

func (f *ResourceFeedAmazonS3) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.S3FeedConfiguration{
		URI:                 resourceDetails["s3_uri"].(string),
		SourceType:          resourceDetails["s3_source_type"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		Authentication: chronicle.S3FeedAuthentication{
			Region:          authenticationDetails["region"].(string),
			AccessKeyID:     authenticationDetails["access_key_id"].(string),
			SecretAccessKey: authenticationDetails["secret_access_key"].(string),
		},
	}
}

//nolint:all
func (f *ResourceFeedAmazonS3) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readS3Conf := readConf.(*chronicle.S3FeedConfiguration)

	// Import Case
	if originalConf == nil {
		return []map[string]interface{}{{
			"s3_uri":                readS3Conf.URI,
			"s3_source_type":        readS3Conf.SourceType,
			"source_delete_options": readS3Conf.SourceDeleteOptions,
			"authentication": []map[string]interface{}{{
				"region":            readS3Conf.Authentication.Region,
				"access_key_id":     readS3Conf.Authentication.AccessKeyID,
				"secret_access_key": readS3Conf.Authentication.SecretAccessKey,
			}},
		}}
	}

	originalS3Conf := originalConf.(*chronicle.S3FeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		"s3_uri":                readS3Conf.URI,
		"s3_source_type":        readS3Conf.SourceType,
		"source_delete_options": originalS3Conf.SourceDeleteOptions, // not returned
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"region":            originalS3Conf.Authentication.Region,
			"access_key_id":     originalS3Conf.Authentication.AccessKeyID,
			"secret_access_key": originalS3Conf.Authentication.SecretAccessKey,
		}},
	}}
}
