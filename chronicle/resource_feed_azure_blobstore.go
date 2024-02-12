package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FeedAzureBlobStoreSourceTypeFiles            = "FILES"
	FeedAzureBlobStoreSourceTypeFolders          = "FOLDERS"
	FeedAzureBlobStoreSourceTypeFoldersRecursive = "FOLDERS_RECURSIVE"
)

const FeedAzureBlobStoreSourceDeleteOptionDeletionNever = "SOURCE_DELETION_NEVER"

type ResourceFeedAzureBlobStore struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedAzureBlobStore() *ResourceFeedAzureBlobStore {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The URI pointing to a Azure Blob Storage blob or container.`,
			},
			"source_type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateFeedAzureBlobStoreSourceType,
				Required:         true,
				Description: `The type of file indicated by the uri. It may be the following:

				- FILES: The URI points to a single file which will be ingested with each execution of the feed.
				- FOLDERS: The URI points to a directory. All files contained within the directory will be ingested with each execution of the feed.
				- FOLDERS_RECURSIVE: The URI points to a directory. All files and directories contains within the indicated directory will be ingested,
				 including all files and directories within those directories, and so on. `,
			},

			"source_delete_options": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateFeedAzureBlobStoreSourceDeleteOption,
				Optional:         true,
				Default:          FeedAzureBlobStoreSourceDeleteOptionDeletionNever,
				Description: `Whether to delete source files after they have been transferred to Chronicle. The possible values are as follows:

				- SOURCE_DELETION_NEVER: Never delete files from the source.`,
			},

			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Azure authentication details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shared_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `A shared key, a 512-bit random string in base64 encoding, authorized to access Azure Blob Storage. Required if not specifying an SAS Token.`,
						},
						"sas_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `A Shared Access Signature authorized to access the Azure Blob Storage container.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a feed from Azure Blobstore service."
	blob := &ResourceFeedAzureBlobStore{}
	blob.TerraformResource = newFeedResourceSchema(details, blob, description, true)

	return blob
}

func (f *ResourceFeedAzureBlobStore) getLogType() string {
	return ""
}

func (f *ResourceFeedAzureBlobStore) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.AzureBlobStoreFeedConfiguration{
		URI:                 resourceDetails["uri"].(string),
		SourceType:          resourceDetails["source_type"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		Authentication: chronicle.AzureBlobStoreFeedFeedAuthentication{
			SharedKey: authenticationDetails["shared_key"].(string),
			SASToken:  authenticationDetails["sas_token"].(string),
		},
	}
}

//nolint:all
func (f *ResourceFeedAzureBlobStore) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {
	readBlobConf := readConf.(*chronicle.AzureBlobStoreFeedConfiguration)

	// Import Case
	if originalConf == nil {
		return []map[string]interface{}{{
			"uri":                   readBlobConf.URI,
			"source_type":           readBlobConf.SourceType,
			"source_delete_options": readBlobConf.SourceDeleteOptions,
			"authentication": []map[string]interface{}{{
				"shared_key": readBlobConf.Authentication.SharedKey,
				"sas_token":  readBlobConf.Authentication.SASToken,
			}},
		}}
	}

	originalBlobConf := originalConf.(*chronicle.AzureBlobStoreFeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		"uri":                   readBlobConf.URI,
		"source_type":           readBlobConf.SourceType,
		"source_delete_options": readBlobConf.SourceDeleteOptions,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"shared_key": originalBlobConf.Authentication.SharedKey,
			"sas_token":  originalBlobConf.Authentication.SASToken,
		}},
	}}
}
