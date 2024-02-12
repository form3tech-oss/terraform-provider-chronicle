package chronicle

import (
	"math/rand"
	"strings"
	"testing"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
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
	result := make([]byte, length)
	set := "abcdefghijklmnopqrstuvwxyz012346789"
	for i := 0; i < length; i++ {
		//nolint:all
		result[i] = set[rand.Intn(len(set))]
	}
	return string(result)
}
