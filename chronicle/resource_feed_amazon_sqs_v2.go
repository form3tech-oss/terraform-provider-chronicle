package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FeedV2SourceDeleteOptionNever     = "NEVER"
	FeedV2SourceDeleteOptionOnSuccess = "ON_SUCCESS"
)

type ResourceFeedAmazonSQSV2 struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedAmazonSQSV2() *ResourceFeedAmazonSQSV2 {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"s3_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The S3 bucket URI in the format s3://bucket-name/path/.`,
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The AWS region where the SQS queue and S3 bucket reside.`,
			},
			"account_number": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validateAWSAccountID,
				Required:         true,
				Description:      `The AWS account number for the SQS queue and S3 bucket.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The SQS queue name.`,
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
				Description: `AWS authentication details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sqs_access_key_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateAWSAccessKeyID,
							Description:      `The 20-character access key ID for the SQS queue.`,
						},
						"sqs_secret_access_key": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validateAWSSecretAccessKey,
							Description:      `The 40-character secret access key for the SQS queue.`,
						},
						"s3_authentication": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `S3 authentication details. Only specify if using a different access key for the S3 bucket.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key_id": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateAWSAccessKeyID,
										Description:      `The 20-character access key ID for the S3 bucket.`,
									},
									"secret_access_key": {
										Type:             schema.TypeString,
										Required:         true,
										Sensitive:        true,
										ValidateDiagFunc: validateAWSSecretAccessKey,
										Description:      `The 40-character secret access key for the S3 bucket.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	description := "Creates a V2 feed from Amazon Simple Queue Service. This feed type uses the Google Cloud Storage Transfer Service for improved ingestion."
	resource := &ResourceFeedAmazonSQSV2{}
	resource.TerraformResource = newFeedResourceSchema(details, resource, description, true)

	return resource
}

func (f *ResourceFeedAmazonSQSV2) getLogType() string {
	return ""
}

func (f *ResourceFeedAmazonSQSV2) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})
	authenticationS3DetailsInterface := authenticationDetails["s3_authentication"].([]interface{})

	sqsConfiguration := &chronicle.SQSV2FeedConfiguration{
		S3URI:               resourceDetails["s3_uri"].(string),
		Region:              resourceDetails["region"].(string),
		AccountNumber:       resourceDetails["account_number"].(string),
		QueueName:           resourceDetails["queue_name"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		MaxLookbackDays:     resourceDetails["max_lookback_days"].(int),
		Authentication: chronicle.SQSV2FeedAuthentication{
			SQSAuthentication: chronicle.SQSV2FeedAuthenticationCred{
				AccessKeyID:     authenticationDetails["sqs_access_key_id"].(string),
				SecretAccessKey: authenticationDetails["sqs_secret_access_key"].(string),
			},
		},
	}

	if len(authenticationS3DetailsInterface) > 0 {
		authenticationS3Details := authenticationS3DetailsInterface[0].(map[string]interface{})
		sqsConfiguration.Authentication.S3Authentication = &chronicle.SQSV2FeedAuthenticationCred{
			AccessKeyID:     authenticationS3Details["access_key_id"].(string),
			SecretAccessKey: authenticationS3Details["secret_access_key"].(string),
		}
	}

	return sqsConfiguration
}

//nolint:all
func (f *ResourceFeedAmazonSQSV2) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readSQSConf := readConf.(*chronicle.SQSV2FeedConfiguration)

	// Import Case
	if originalConf == nil {
		conf := []map[string]interface{}{{
			"s3_uri":                readSQSConf.S3URI,
			"region":                readSQSConf.Region,
			"account_number":        readSQSConf.AccountNumber,
			"queue_name":            readSQSConf.QueueName,
			"source_delete_options": readSQSConf.SourceDeleteOptions,
			"max_lookback_days":     readSQSConf.MaxLookbackDays,
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

	originalSQSConf := originalConf.(*chronicle.SQSV2FeedConfiguration)
	// Default Case
	conf := []map[string]interface{}{{
		"s3_uri":                readSQSConf.S3URI,
		"region":                readSQSConf.Region,
		"account_number":        readSQSConf.AccountNumber,
		"queue_name":            readSQSConf.QueueName,
		"source_delete_options": originalSQSConf.SourceDeleteOptions, // not returned
		"max_lookback_days":     readSQSConf.MaxLookbackDays,
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
