package chronicle

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readStringFromResource(d *schema.ResourceData, key string) string {
	if attr, ok := d.GetOk(key); ok {
		return attr.(string)
	}
	return ""
}

func readBoolFromResource(d *schema.ResourceData, key string) bool {
	if attr, ok := d.GetOk(key); ok {
		return attr.(bool)
	}
	return false
}

func readMapFromResource(d *schema.ResourceData, key string) map[string]interface{} {
	if attr, ok := d.GetOk(key); ok {
		result := attr.(map[string]interface{})
		return result
	}

	return nil
}

func readStringSliceFromResource(d *schema.ResourceData, key string) []string {
	if attr, ok := d.GetOk(key); ok {
		var array []string
		items, ok := attr.([]interface{})
		if !ok {
			return nil
		}
		for _, x := range items {
			array = append(array, x.(string))
		}

		return array
	}

	return nil
}

//nolint:unparam
func readSliceFromResource(d *schema.ResourceData, key string) []interface{} {
	if attr, ok := d.GetOk(key); ok {
		var slice []interface{}
		items := attr.([]interface{})
		slice = append(slice, items...)

		return slice
	}

	return nil
}
