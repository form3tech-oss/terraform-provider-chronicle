package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedOktaSystemLog struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedOktaSystemLog() *ResourceFeedOktaSystemLog {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Okta hostname.`,
			},
			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Okta authentication header details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Okta authorization key.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `Okta API token.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a feed from API source Type for Okta System Log log type."
	okta := &ResourceFeedOktaSystemLog{}
	okta.TerraformResource = newFeedResourceSchema(details, okta, description, false)

	return okta
}

func (f *ResourceFeedOktaSystemLog) getLogType() string {
	return chronicle.OktaSystemLogFeedLogType
}

func (f *ResourceFeedOktaSystemLog) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.OktaSystemLogFeedConfiguration{
		Hostname: resourceDetails["hostname"].(string),
		Authentication: chronicle.OktaSystemLogFeedAuthentication{
			HeaderKeyValues: []chronicle.OktaSystemLogAuthenticationHeaderKeyValues{
				{
					Key:   authenticationDetails["key"].(string),
					Value: authenticationDetails["value"].(string),
				},
			},
		},
	}
}

//nolint:all
func (f *ResourceFeedOktaSystemLog) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readOktaConf := readConf.(*chronicle.OktaSystemLogFeedConfiguration)

	// Import Case
	if originalConf == nil {
		var key, value string
		if readOktaConf.Authentication.HeaderKeyValues == nil {
			key = ""
			value = ""
		} else {
			key = readOktaConf.Authentication.HeaderKeyValues[0].Key
			value = readOktaConf.Authentication.HeaderKeyValues[0].Value
		}

		return []map[string]interface{}{{
			"hostname": readOktaConf.Hostname,
			"authentication": []map[string]interface{}{{
				"key":   key,
				"value": value,
			},
			}},
		}
	}

	originalOkta := originalConf.(*chronicle.OktaSystemLogFeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		"hostname": readOktaConf.Hostname,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"key":   originalOkta.Authentication.HeaderKeyValues[0].Key,
			"value": originalOkta.Authentication.HeaderKeyValues[0].Value,
		},
		}},
	}
}
