package chronicle

import (
	"strings"
	"testing"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"chronicle": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	t.Parallel()
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := multiEnvSearch(chronicle.EnvAPICrendetialsVar); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(chronicle.EnvAPICrendetialsVar, ", "))
	}
}

func randString(length int) string {
	id := uuid.New().String()

	if len(id) > length {
		id = id[:length]
	}

	return id
}
