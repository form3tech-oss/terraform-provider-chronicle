package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FeedMicrosoftOffice365ManagementActivityContentTypeAuditAzureActiveDirectory = "AUDIT_AZURE_ACTIVE_DIRECTORY"
	FeedMicrosoftOffice365ManagementActivityContentTypeAuditExchange             = "AUDIT_EXCHANGE"
	FeedMicrosoftOffice365ManagementActivityContentTypeAuditSharePoint           = "AUDIT_SHARE_POINT"
	FeedMicrosoftOffice365ManagementActivityContentTypeAuditGeneral              = "AUDIT_GENERAL"
	FeedMicrosoftOffice365ManagementActivityContentTypeDPLAll                    = "DLP_ALL"
)

type ResourceFeedMicrosoftOffice365ManagementActivity struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedMicrosoftOffice365ManagementActivity() *ResourceFeedMicrosoftOffice365ManagementActivity {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `API Full Path, default value: manage.office.com/api/v1.0.`,
			},
			"tenant_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateUUID,
				Description:      `Tenant ID (a UUID).`,
			},
			"content_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateFeedMicrosoftOffice365ManagementActivityContentType,
				Description:      `The type of logs to fetch. See https://cloud.google.com/chronicle/docs/reference/feed-management-api#office_365_content_type.`,
			},
			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Office 365 authentication details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateUUID,
							Description:      `OAuth client ID (a UUID).`,
						},
						"client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `OAuth client secret.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a feed from API source type for Microsoft Office 365 Management Activity log type."
	microsoftOffice365ManagementActivity := &ResourceFeedMicrosoftOffice365ManagementActivity{}
	microsoftOffice365ManagementActivity.TerraformResource = newFeedResourceSchema(details, microsoftOffice365ManagementActivity, description, false)

	return microsoftOffice365ManagementActivity
}

func (f *ResourceFeedMicrosoftOffice365ManagementActivity) getLogType() string {
	return chronicle.MicrosoftOffice365ManagementActivityFeedLogType
}

func (f *ResourceFeedMicrosoftOffice365ManagementActivity) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.MicrosoftOffice365ManagementActivityFeedConfiguration{
		Hostname:    resourceDetails["hostname"].(string),
		TenantID:    resourceDetails["tenant_id"].(string),
		ContentType: resourceDetails["content_type"].(string),
		Authentication: chronicle.MicrosoftOffice365ManagementActivityFeedAuthentication{
			ClientID:     authenticationDetails["client_id"].(string),
			ClientSecret: authenticationDetails["client_secret"].(string),
		},
	}
}

//nolint:all
func (f *ResourceFeedMicrosoftOffice365ManagementActivity) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readMicrosoftOffice365ManagementActivityConf := readConf.(*chronicle.MicrosoftOffice365ManagementActivityFeedConfiguration)

	// Import Case
	if originalConf == nil {
		return []map[string]interface{}{{
			"hostname":     readMicrosoftOffice365ManagementActivityConf.Hostname,
			"tenant_id":    readMicrosoftOffice365ManagementActivityConf.TenantID,
			"content_type": readMicrosoftOffice365ManagementActivityConf.ContentType,
			"authentication": []map[string]interface{}{{
				"client_id":     readMicrosoftOffice365ManagementActivityConf.Authentication.ClientID,
				"client_secret": readMicrosoftOffice365ManagementActivityConf.Authentication.ClientSecret,
			}},
		}}
	}

	originalMicrosoftOffice365ManagementActivity := originalConf.(*chronicle.MicrosoftOffice365ManagementActivityFeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		// reponse strips out API endpoints. Workaround: get it from original
		"hostname":     originalMicrosoftOffice365ManagementActivity.Hostname,
		"tenant_id":    readMicrosoftOffice365ManagementActivityConf.TenantID,
		"content_type": readMicrosoftOffice365ManagementActivityConf.ContentType,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"client_id":     originalMicrosoftOffice365ManagementActivity.Authentication.ClientID,
			"client_secret": originalMicrosoftOffice365ManagementActivity.Authentication.ClientSecret,
		}},
	}}
}
