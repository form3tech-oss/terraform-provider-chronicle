package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceOktaUsers struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedOktaUsers() *ResourceOktaUsers {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Okta hostname.`,
			},
			"manager_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Manager ID is required when you use a non-Okta ID to reference managers.
				It should be a JSON field path pointing to the field that contains the manager ID in the result of a call to the "users" Okta API.`,
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
	description := "Creates a feed from API source Type for Okta Users log type."
	okta := &ResourceOktaUsers{}
	okta.TerraformResource = newFeedResourceSchema(details, okta, description, false)

	return okta
}

func (f *ResourceOktaUsers) getLogType() string {
	return chronicle.OktaUsersFeedLogType
}

func (f *ResourceOktaUsers) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.OktaUsersFeedConfiguration{
		Hostname:                resourceDetails["hostname"].(string),
		ManagerIDReferenceField: resourceDetails["manager_id"].(string),
		Authentication: chronicle.OktaUsersFeedAuthentication{
			HeaderKeyValues: []chronicle.OktaUsersAuthenticationHeaderKeyValues{
				{
					Key:   authenticationDetails["key"].(string),
					Value: authenticationDetails["value"].(string),
				},
			},
		},
	}
}

//nolint:all
func (f *ResourceOktaUsers) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readOktaConf := readConf.(*chronicle.OktaUsersFeedConfiguration)

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
			"hostname":   readOktaConf.Hostname,
			"manager_id": readOktaConf.ManagerIDReferenceField,
			"authentication": []map[string]interface{}{{
				"key":   key,
				"value": value,
			},
			}},
		}
	}

	originalOkta := originalConf.(*chronicle.OktaUsersFeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		"hostname":   originalOkta.Hostname,
		"manager_id": originalOkta.ManagerIDReferenceField,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"key":   originalOkta.Authentication.HeaderKeyValues[0].Key,
			"value": originalOkta.Authentication.HeaderKeyValues[0].Value,
		},
		}},
	}
}
