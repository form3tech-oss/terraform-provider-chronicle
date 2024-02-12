package chronicle

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"
)

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func envSearch(s string) string {
	if v := os.Getenv(s); v != "" {
		return v
	}

	return ""
}

// checks if a string is present in a slice.
func contains(s []string, s1 string) bool {
	for _, a := range s {
		if a == s1 {
			return true
		}
	}
	return false
}

func handleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if isGoogleAPIErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore.
		d.SetId("")

		return nil
	}

	return fmt.Errorf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resource), err)
}
func isGoogleAPIErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	return ok && gerr != nil && gerr.Code == errCode
}
