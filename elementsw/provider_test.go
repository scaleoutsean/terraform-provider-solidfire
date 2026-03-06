package elementsw

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
}

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"solidfire":       func() (*schema.Provider, error) { return testAccProvider, nil },
	"solidfireremote": func() (*schema.Provider, error) { return testAccProvider, nil },
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SOLIDFIRE_USERNAME"); v == "" {
		t.Fatal("SOLIDFIRE_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("SOLIDFIRE_PASSWORD"); v == "" {
		t.Fatal("SOLIDFIRE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("SOLIDFIRE_SERVER"); v == "" {
		t.Fatal("SOLIDFIRE_SERVER must be set for acceptance tests")
	}

	if v := os.Getenv("SOLIDFIRE_API_VERSION"); v == "" {
		t.Fatal("SOLIDFIRE_API_VERSION must be set for acceptance tests")
	}
}
