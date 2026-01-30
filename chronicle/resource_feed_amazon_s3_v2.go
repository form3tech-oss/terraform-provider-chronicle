package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedAmazonS3V2 struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedAmazonS3V2() *ResourceFeedAmazonS3V2 {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
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
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `AWS authentication details. Use either access key credentials or IAM role ARN.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateAWSAccessKeyID,
							ConflictsWith:    []string{"details.0.authentication.0.aws_iam_role_arn"},
							Description:      `The 20-character access key ID associated with your Amazon IAM account. Required if not using aws_iam_role_arn.`,
						},
						"secret_access_key": {
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateDiagFunc: validateAWSSecretAccessKey,
							ConflictsWith:    []string{"details.0.authentication.0.aws_iam_role_arn"},
							Description:      `The 40-character secret access key associated with your Amazon IAM account. Required if not using aws_iam_role_arn.`,
						},
						"aws_iam_role_arn": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"details.0.authentication.0.access_key_id", "details.0.authentication.0.secret_access_key"},
							Description:   `ARN of the AWS IAM role configured to access S3 bucket. Use this for federated authentication instead of access keys.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a V2 feed from Amazon Simple Storage Service (S3). " +
		"This feed type uses the Google Cloud Storage Transfer Service for improved ingestion."
	resource := &ResourceFeedAmazonS3V2{}
	resource.TerraformResource = newFeedResourceSchema(details, resource, description, true)

	return resource
}

func (f *ResourceFeedAmazonS3V2) getLogType() string {
	return ""
}

func (f *ResourceFeedAmazonS3V2) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	config := &chronicle.S3V2FeedConfiguration{
		S3URI:               resourceDetails["s3_uri"].(string),
		SourceDeleteOptions: resourceDetails["source_delete_options"].(string),
		MaxLookbackDays:     resourceDetails["max_lookback_days"].(int),
		Authentication:      chronicle.S3V2FeedAuthentication{},
	}

	// Check which authentication method is used
	if iamRoleArn, ok := authenticationDetails["aws_iam_role_arn"].(string); ok && iamRoleArn != "" {
		config.Authentication.AWSIAMRoleArn = iamRoleArn
	} else {
		config.Authentication.AccessKeySecretAuth = &chronicle.S3V2AccessKeySecretAuth{
			AccessKeyID:     authenticationDetails["access_key_id"].(string),
			SecretAccessKey: authenticationDetails["secret_access_key"].(string),
		}
	}

	return config
}

//nolint:all
func (f *ResourceFeedAmazonS3V2) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readS3Conf := readConf.(*chronicle.S3V2FeedConfiguration)

	// Import Case
	if originalConf == nil {
		authMap := make(map[string]interface{})
		if readS3Conf.Authentication.AWSIAMRoleArn != "" {
			authMap["aws_iam_role_arn"] = readS3Conf.Authentication.AWSIAMRoleArn
		}
		if readS3Conf.Authentication.AccessKeySecretAuth != nil {
			authMap["access_key_id"] = readS3Conf.Authentication.AccessKeySecretAuth.AccessKeyID
			authMap["secret_access_key"] = readS3Conf.Authentication.AccessKeySecretAuth.SecretAccessKey
		}

		return []map[string]interface{}{{
			"s3_uri":                readS3Conf.S3URI,
			"source_delete_options": readS3Conf.SourceDeleteOptions,
			"max_lookback_days":     readS3Conf.MaxLookbackDays,
			"authentication":        []map[string]interface{}{authMap},
		}}
	}

	originalS3Conf := originalConf.(*chronicle.S3V2FeedConfiguration)
	// Default Case
	authMap := make(map[string]interface{})
	if originalS3Conf.Authentication.AWSIAMRoleArn != "" {
		authMap["aws_iam_role_arn"] = originalS3Conf.Authentication.AWSIAMRoleArn
	}
	if originalS3Conf.Authentication.AccessKeySecretAuth != nil {
		authMap["access_key_id"] = originalS3Conf.Authentication.AccessKeySecretAuth.AccessKeyID
		authMap["secret_access_key"] = originalS3Conf.Authentication.AccessKeySecretAuth.SecretAccessKey
	}

	return []map[string]interface{}{{
		"s3_uri":                readS3Conf.S3URI,
		"source_delete_options": originalS3Conf.SourceDeleteOptions, // not returned
		"max_lookback_days":     readS3Conf.MaxLookbackDays,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{authMap},
	}}
}
