package chronicle

import (
	"fmt"
	"log"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewNotFoundErrorf(format string, a ...interface{}) error {
	return fmt.Errorf("%s %s", "Could not find", fmt.Sprintf(format, a...))
}

func HandleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if IsChronicleAPIErrorWithCode(err, 404) {
		log.Printf("Removing %s because it's gone", resource)
		d.SetId("")

		return nil
	}

	return fmt.Errorf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resource), err)
}

func IsChronicleAPIErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &chronicle.ChronicleAPIError{}).(*chronicle.ChronicleAPIError)
	return ok && gerr != nil && gerr.HTTPStatusCode == errCode
}
