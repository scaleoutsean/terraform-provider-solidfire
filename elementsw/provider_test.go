package elementsw

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"elementsw":       func() (*schema.Provider, error) { return Provider(), nil },
	"elementswremote": func() (*schema.Provider, error) { return Provider(), nil },
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func init() {
	// No global testAccProvider to avoid leakage between tests
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ELEMENTSW_USERNAME"); v == "" {
		t.Fatal("ELEMENTSW_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("ELEMENTSW_PASSWORD"); v == "" {
		t.Fatal("ELEMENTSW_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("ELEMENTSW_SERVER"); v == "" {
		t.Fatal("ELEMENTSW_SERVER must be set for acceptance tests")
	}

	if v := os.Getenv("ELEMENTSW_API_VERSION"); v == "" {
		t.Fatal("ELEMENTSW_API_VERSION must be set for acceptance tests")
	}
}
