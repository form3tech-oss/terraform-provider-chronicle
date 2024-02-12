package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedThinkstCanary struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedThinkstCanary() *ResourceFeedThinkstCanary {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateThinkstCanaryHostname,
				Description:      `Thinkst Canary hostname.`,
			},
			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Thinkst Canary authentication details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "auth_token",
							Description: `Thinkst Canary authentication key. Defaults to auth_token.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `Thinkst Canary authentication value.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a feed from API source Type for Thinkst Canary log type."
	canary := &ResourceFeedThinkstCanary{}
	canary.TerraformResource = newFeedResourceSchema(details, canary, description, false)

	return canary
}

func (f *ResourceFeedThinkstCanary) getLogType() string {
	return chronicle.ThinkstCanaryFeedLogType
}

func (f *ResourceFeedThinkstCanary) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.ThinkstCanaryFeedConfiguration{
		Hostname: resourceDetails["hostname"].(string),
		Authentication: chronicle.ThinkstCanaryFeedAuthentication{
			HeaderKeyValues: []chronicle.ThinkstCanaryAuthenticationHeaderKeyValues{
				{
					Key:   authenticationDetails["key"].(string),
					Value: authenticationDetails["value"].(string),
				},
			},
		},
	}
}

//nolint:all
func (f *ResourceFeedThinkstCanary) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readCanaryConf := readConf.(*chronicle.ThinkstCanaryFeedConfiguration)

	// Import Case
	if originalConf == nil {
		var key, value string
		if readCanaryConf.Authentication.HeaderKeyValues == nil {
			key = ""
			value = ""
		} else {
			key = readCanaryConf.Authentication.HeaderKeyValues[0].Key
			value = readCanaryConf.Authentication.HeaderKeyValues[0].Value
		}

		return []map[string]interface{}{{
			"hostname": readCanaryConf.Hostname,
			"authentication": []map[string]interface{}{{
				"key":   key,
				"value": value,
			},
			}},
		}
	}

	originalCanary := originalConf.(*chronicle.ThinkstCanaryFeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		"hostname": readCanaryConf.Hostname,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"key":   originalCanary.Authentication.HeaderKeyValues[0].Key,
			"value": originalCanary.Authentication.HeaderKeyValues[0].Value,
		},
		}},
	}
}
