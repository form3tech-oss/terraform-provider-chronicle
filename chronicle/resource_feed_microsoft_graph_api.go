package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FeedMicrosoftGraphAPIContentTypeAzureADAudit   = "AZURE_AD_AUDIT"
	FeedMicrosoftGraphAPIContentTypeAzureADContext = "AZURE_AD_CONTEXT"
	FeedMicrosoftGraphAPIContentTypeAzureSignIns   = "AZURE_AD"
	FeedMicrosoftGraphAPIContentTypeAzureMDMIntune = "AZURE_MDM_INTUNE"
	FeedMicrosoftGraphAPIContentTypeMSGraphAlert   = "MICROSOFT_GRAPH_ALERT"
)

type ResourceFeedMicrosoftGraphAPI struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedMicrosoftGraphAPI() *ResourceFeedMicrosoftGraphAPI {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `API Full Path, default value: e.g. graph.microsoft.com/beta/security/alerts_v2`,
			},
			"tenant_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateUUID,
				Description:      `Tenant ID (a UUID).`,
			},
			"retrieve_devices": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to retrieve devices from the directory when using AZURE_AD_CONTEXT.`,
			},
			"retrieve_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to retrieve groups from the directory when using AZURE_AD_CONTEXT.`,
			},
			"auth_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Authentication endpoint to use when using MICROSOFT_GRAPH_ALERT.`,
			},
			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Microsoft Graph API authentication details.`,
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
	description := "Creates a feed from API source type for a Microsoft Graph API log type."
	microsoftGraphAPI := &ResourceFeedMicrosoftGraphAPI{}
	microsoftGraphAPI.TerraformResource = newFeedResourceSchema(details, microsoftGraphAPI, description, true)

	return microsoftGraphAPI
}

func (f *ResourceFeedMicrosoftGraphAPI) getLogType() string {
	return ""
}

func (f *ResourceFeedMicrosoftGraphAPI) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {

	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	contentType := readStringFromResource(d, "log_type")

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.MicrosoftGraphAPIFeedConfiguration{
		Hostname:        resourceDetails["hostname"].(string),
		TenantID:        resourceDetails["tenant_id"].(string),
		RetrieveDevices: resourceDetails["retrieve_devices"].(bool),
		RetrieveGroups:  resourceDetails["retrieve_groups"].(bool),
		AuthEndpoint:    resourceDetails["auth_endpoint"].(string),
		ContentType:     contentType,
		Authentication: chronicle.MicrosoftGraphAPIFeedAuthentication{
			ClientID:     authenticationDetails["client_id"].(string),
			ClientSecret: authenticationDetails["client_secret"].(string),
		},
	}
}

//nolint:all
func (f *ResourceFeedMicrosoftGraphAPI) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {
	readMicrosoftGraphAPIConf := readConf.(*chronicle.MicrosoftGraphAPIFeedConfiguration)
	// Import Case
	if originalConf == nil {
		return []map[string]interface{}{{
			"hostname":         readMicrosoftGraphAPIConf.Hostname,
			"tenant_id":        readMicrosoftGraphAPIConf.TenantID,
			"retrieve_devices": readMicrosoftGraphAPIConf.RetrieveDevices,
			"retrieve_groups":  readMicrosoftGraphAPIConf.RetrieveGroups,
			"auth_endpoint":    readMicrosoftGraphAPIConf.AuthEndpoint,
			"authentication": []map[string]interface{}{{
				"client_id":     readMicrosoftGraphAPIConf.Authentication.ClientID,
				"client_secret": readMicrosoftGraphAPIConf.Authentication.ClientSecret,
			}},
		}}
	}

	originalMicrosoftGraphAPI := originalConf.(*chronicle.MicrosoftGraphAPIFeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		// reponse strips out API endpoints. Workaround: get it from original
		"hostname":         originalMicrosoftGraphAPI.Hostname,
		"tenant_id":        readMicrosoftGraphAPIConf.TenantID,
		"retrieve_devices": readMicrosoftGraphAPIConf.RetrieveDevices,
		"retrieve_groups":  readMicrosoftGraphAPIConf.RetrieveGroups,
		"auth_endpoint":    readMicrosoftGraphAPIConf.AuthEndpoint,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"client_id":     originalMicrosoftGraphAPI.Authentication.ClientID,
			"client_secret": originalMicrosoftGraphAPI.Authentication.ClientSecret,
		}},
	}}
}
