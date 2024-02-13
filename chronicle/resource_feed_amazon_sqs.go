package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedAmazonSQS struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedAmazonSQS() *ResourceFeedAmazonSQS {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"queue": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The SQS queue name.`,
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The region where the S3 bucket resides following this format: https://cloud.google.com/chronicle/docs/reference/feed-management-api#amazon_s3_regions.`,
			},

			"account_number": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateAWSAccountID,
				Required:         true,
				Description:      `The account number for the SQS queue and S3 bucket.`,
			},

			"source_delete_options": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateFeedS3SourceDeleteOption,
				Required:         true,
				Description: `Whether to delete the source files in the S3 bucket after they have been transferred to Chronicle. This reduces storage costs. Valid values are:

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
						"sqs_access_key_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateAWSAccessKeyID,
							Description:      `This is the 20 character ID associated with your Amazon IAM account.`,
						},
						"sqs_secret_access_key": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validateAWSSecretAccessKey,
							Description:      `This is the 40 character access key associated with your Amazon IAM account.`,
						},
						"s3_authentication": {
							Type:        schema.TypeList,
							Required:    false,
							Optional:    true,
							MaxItems:    1,
							Description: `S3 AWS authentication details. Only specify if using a different access key for the S3 bucket.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
	description := "Creates a feed from Amazon Simple Queue Service."
	sqs := &ResourceFeedAmazonSQS{}
	sqs.TerraformResource = newFeedResourceSchema(details, sqs, description, true)

	return sqs
}

func (f *ResourceFeedAmazonSQS) getLogType() string {
	return ""
}

func (f *ResourceFeedAmazonSQS) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})
	authenticationS3DetailsInterface := authenticationDetails["s3_authentication"].([]interface{})

	sqsConfiguration := &chronicle.SQSFeedConfiguration{
		Queue:               resourceDetails["queue"].(string),
		Region:              resourceDetails["region"].(string),
		AccountNumber:       resourceDetails["account_number"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		Authentication: chronicle.SQSFeedAuthentication{
			SQSAuthentication: chronicle.SQSFeedAuthenticationCred{
				AccessKeyID:     authenticationDetails["sqs_access_key_id"].(string),
				SecretAccessKey: authenticationDetails["sqs_secret_access_key"].(string),
			},
		},
	}

	if len(authenticationS3DetailsInterface) > 0 {
		authenticationS3Details := authenticationS3DetailsInterface[0].(map[string]interface{})
		sqsConfiguration.Authentication.S3Authentication = &chronicle.SQSFeedAuthenticationCred{
			AccessKeyID:     authenticationS3Details["access_key_id"].(string),
			SecretAccessKey: authenticationS3Details["secret_access_key"].(string),
		}
	}

	return sqsConfiguration
}

//nolint:all
func (f *ResourceFeedAmazonSQS) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readSQSConf := readConf.(*chronicle.SQSFeedConfiguration)

	// Import Case
	if originalConf == nil {
		conf := []map[string]interface{}{{
			"queue":                 readSQSConf.Queue,
			"region":                readSQSConf.Region,
			"account_number":        readSQSConf.AccountNumber,
			"source_delete_options": readSQSConf.SourceDeleteOptions,
			"authentication": []map[string]interface{}{{
				"sqs_access_key_id":     readSQSConf.Authentication.SQSAuthentication.AccessKeyID,
				"sqs_secret_access_key": readSQSConf.Authentication.SQSAuthentication.SecretAccessKey,
			}},
		}}

		if readSQSConf.Authentication.S3Authentication != nil {
			conf[0]["authentication"].([]map[string]interface{})[0]["s3_authentication"] = []map[string]interface{}{{
				"access_key_id":     readSQSConf.Authentication.S3Authentication.AccessKeyID,
				"secret_access_key": readSQSConf.Authentication.S3Authentication.SecretAccessKey,
			}}
		}

		return conf
	}

	originalSQSConf := originalConf.(*chronicle.SQSFeedConfiguration)
	// Default Case
	conf := []map[string]interface{}{{
		"queue":                 readSQSConf.Queue,
		"region":                readSQSConf.Region,
		"account_number":        readSQSConf.AccountNumber,
		"source_delete_options": originalSQSConf.SourceDeleteOptions, // not returned
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"sqs_access_key_id":     originalSQSConf.Authentication.SQSAuthentication.AccessKeyID,
			"sqs_secret_access_key": originalSQSConf.Authentication.SQSAuthentication.SecretAccessKey,
		}},
	}}

	if originalSQSConf.Authentication.S3Authentication != nil {
		conf[0]["authentication"].([]map[string]interface{})[0]["s3_authentication"] = []map[string]interface{}{{
			"access_key_id":     originalSQSConf.Authentication.S3Authentication.AccessKeyID,
			"secret_access_key": originalSQSConf.Authentication.S3Authentication.SecretAccessKey,
		}}
	}

	return conf
}
