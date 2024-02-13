package chronicle

import (
	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceFeedQualysVM struct {
	TerraformResource *schema.Resource
}

func NewResourceFeedQualysVM() *ResourceFeedQualysVM {
	details := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Qualys VM hostname.`,
			},
			"authentication": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `AWS authentication details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Username.`,
						},
						"secret": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `Password.`,
						},
					},
				},
			},
		},
	}
	description := "Creates a feed from API source Type for Qualys VM log type."
	qualysVM := &ResourceFeedQualysVM{}
	qualysVM.TerraformResource = newFeedResourceSchema(details, qualysVM, description, false)

	return qualysVM
}

func (f *ResourceFeedQualysVM) getLogType() string {
	return chronicle.QualysVMFeedLogType
}

func (f *ResourceFeedQualysVM) expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration {
	resourceDetailsInterface := readSliceFromResource(d, "details")
	if resourceDetailsInterface == nil {
		return nil
	}

	resourceDetails := resourceDetailsInterface[0].(map[string]interface{})
	authenticationDetails := resourceDetails["authentication"].([]interface{})[0].(map[string]interface{})

	return &chronicle.QualysVMFeedConfiguration{
		Hostname: resourceDetails["hostname"].(string),
		Authentication: chronicle.QualysVMFeedAuthentication{
			User:   authenticationDetails["user"].(string),
			Secret: authenticationDetails["secret"].(string),
		},
	}
}

//nolint:all
func (f *ResourceFeedQualysVM) flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{} {

	readQualysVMConf := readConf.(*chronicle.QualysVMFeedConfiguration)

	// Import Case
	if originalConf == nil {
		return []map[string]interface{}{{
			"hostname": readQualysVMConf.Hostname,
			"authentication": []map[string]interface{}{{
				"user":   readQualysVMConf.Authentication.User,
				"secret": readQualysVMConf.Authentication.Secret,
			}},
		}}
	}

	originalQualysVM := originalConf.(*chronicle.QualysVMFeedConfiguration)
	// Default Case
	return []map[string]interface{}{{
		"hostname": originalQualysVM.Hostname,
		// replace authentication block with original values because they are not returned within a read request
		"authentication": []map[string]interface{}{{
			"user":   originalQualysVM.Authentication.User,
			"secret": originalQualysVM.Authentication.Secret,
		}},
	}}
}
