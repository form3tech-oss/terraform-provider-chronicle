package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedProofpointSIEM struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedProofpointSIEM() *ResourceFeedProofpointSIEM {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Proofpoint authentication header details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Proofpoint user.`,
						},
						"secret": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `Proofpoint secret.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a feed from API source type for Proofpoint SIEM log type."
	proofpoint := &ResourceFeedProofpointSIEM{}
	proofpoint.TerraformResource = newFeedResourceSchema(details, proofpoint, description, false)

	return proofpoint
}

func (f *ResourceFeedProofpointSIEM) getLogType() string {
	return chronicle.ProofpointSIEMFeedLogType
}

func (f *ResourceFeedProofpointSIEM) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.ProofpointSIEMFeedConfiguration{
		Authentication: chronicle.ProofpointSIEMFeedAuthentication{
			User:   authenticationDetails["user"].(string),
			Secret: authenticationDetails["secret"].(string),
		},
	}
}

//nolint:all
func (f *ResourceFeedProofpointSIEM) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {
	readProofpointConf := readConf.(*chronicle.ProofpointSIEMFeedConfiguration)
	if originalConf == nil {
		return []map[string]interface{}{{
			"authentication": []map[string]interface{}{{
				"user":   readProofpointConf.Authentication.User,
				"secret": readProofpointConf.Authentication.Secret,
			},
			}},
		}
	}

	originalProofpointConf := originalConf.(*chronicle.ProofpointSIEMFeedConfiguration)
	return []map[string]interface{}{{
		"authentication": []map[string]interface{}{{
			"user":   originalProofpointConf.Authentication.User,
			"secret": originalProofpointConf.Authentication.Secret,
		},
		}},
	}

}
