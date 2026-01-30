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
			"queue": {
				Type:     schema.TypeString,
				Required: true,
				Description: `The Amazon Resource Name (ARN) of the SQS queue. ` +
					`Format: arn:aws:sqs:region:account_id:queue_name. Example: arn:aws:sqs:us-east-1:123456789012:my-queue`,
			},
			"s3_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The S3 bucket URI in the format s3://bucket-name/path/.`,
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
			"authentication": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Description: `AWS authentication details. Use either access key credentials or IAM role ARN. ` +
					`The same credentials are used for both SQS queue and S3 bucket.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateAWSAccessKeyID,
							ConflictsWith:    []string{"details.0.authentication.0.aws_iam_role_arn"},
							Description: `The 20-character access key ID for your Amazon IAM account. ` +
								`Required if not using aws_iam_role_arn. Same credentials are used for both SQS queue and S3 bucket.`,
						},
						"secret_access_key": {
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateDiagFunc: validateAWSSecretAccessKey,
							ConflictsWith:    []string{"details.0.authentication.0.aws_iam_role_arn"},
							Description: `The 40-character secret access key for your Amazon IAM account. ` +
								`Required if not using aws_iam_role_arn. Same credentials are used for both SQS queue and S3 bucket.`,
						},
						"aws_iam_role_arn": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"details.0.authentication.0.access_key_id", "details.0.authentication.0.secret_access_key"},
							Description: `ARN of the AWS IAM role configured to access both SQS queue and S3 bucket. ` +
								`Use this for federated authentication instead of access keys.`,
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

	config := &chronicle.SQSV2FeedConfiguration{
		Queue:               resourceDetails["queue"].(string),
		S3URI:               resourceDetails["s3_uri"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		MaxLookbackDays:     resourceDetails["max_lookback_days"].(int),
		Authentication:      chronicle.SQSV2FeedAuthentication{},
	}

	// Check which authentication method is used
	if iamRoleArn, ok := authenticationDetails["aws_iam_role_arn"].(string); ok && iamRoleArn != "" {
		config.Authentication.AWSIAMRoleArn = iamRoleArn
	} else {
		config.Authentication.AccessKeySecretAuth = &chronicle.SQSV2AccessKeySecretAuth{
			AccessKeyID:     authenticationDetails["access_key_id"].(string),
			SecretAccessKey: authenticationDetails["secret_access_key"].(string),
		}
	}

	return config
}

//nolint:all
func (f *ResourceFeedAmazonSQSV2) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readSQSConf := readConf.(*chronicle.SQSV2FeedConfiguration)

	// Import Case
	if originalConf == nil {
		authMap := make(map[string]interface{})
		if readSQSConf.Authentication.AWSIAMRoleArn != "" {
			authMap["aws_iam_role_arn"] = readSQSConf.Authentication.AWSIAMRoleArn
		}
		if readSQSConf.Authentication.AccessKeySecretAuth != nil {
			authMap["access_key_id"] = readSQSConf.Authentication.AccessKeySecretAuth.AccessKeyID
			authMap["secret_access_key"] = readSQSConf.Authentication.AccessKeySecretAuth.SecretAccessKey
		}

		return []map[string]interface{}{{
			"queue":                 readSQSConf.Queue,
			"s3_uri":                readSQSConf.S3URI,
			"source_delete_options": readSQSConf.SourceDeleteOptions,
			"max_lookback_days":     readSQSConf.MaxLookbackDays,
			"authentication":        []map[string]interface{}{authMap},
		}}
	}

	originalSQSConf := originalConf.(*chronicle.SQSV2FeedConfiguration)
	// Default Case
	authMap := make(map[string]interface{})
	if originalSQSConf.Authentication.AWSIAMRoleArn != "" {
		authMap["aws_iam_role_arn"] = originalSQSConf.Authentication.AWSIAMRoleArn
	}
	if originalSQSConf.Authentication.AccessKeySecretAuth != nil {
		authMap["access_key_id"] = originalSQSConf.Authentication.AccessKeySecretAuth.AccessKeyID
		authMap["secret_access_key"] = originalSQSConf.Authentication.AccessKeySecretAuth.SecretAccessKey
	}

	return []map[string]interface{}{{
		"queue":                 readSQSConf.Queue,
		"s3_uri":                readSQSConf.S3URI,
		"source_delete_options": originalSQSConf.SourceDeleteOptions, // not returned
		"max_lookback_days":     readSQSConf.MaxLookbackDays,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{authMap},
	}}
}
