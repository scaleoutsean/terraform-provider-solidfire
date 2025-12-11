package elementsw

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccElementswReplication_basic(t *testing.T) {
	drServer := os.Getenv("ELEMENTSW_SERVER_DR")
	if drServer == "" {
		t.Skip("ELEMENTSW_SERVER_DR not set, skipping replication tests")
	}

	// Get primary credentials (required for source_cluster in automated pairing)
	srcServer := os.Getenv("ELEMENTSW_SERVER")
	srcUser := os.Getenv("ELEMENTSW_USERNAME")
	srcPass := os.Getenv("ELEMENTSW_PASSWORD")

	// Get DR credentials (default to primary if not set)
	drUser := os.Getenv("ELEMENTSW_USERNAME_DR")
	if drUser == "" {
		drUser = srcUser
	}
	drPass := os.Getenv("ELEMENTSW_PASSWORD_DR")
	if drPass == "" {
		drPass = srcPass
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReplicationClusterConfig(srcServer, srcUser, srcPass, drServer, drUser, drPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("elementsw_replication_cluster.test", "cluster_pair_id"),
					resource.TestCheckResourceAttrSet("elementsw_replication_cluster.test", "cluster_name"),
				),
			},
		},
	})
}

func testAccReplicationClusterConfig(srcServer, srcUser, srcPass, drServer, drUser, drPass string) string {
	return fmt.Sprintf(`
resource "elementsw_replication_cluster" "test" {
  source_cluster {
    endpoint = "%s"
    username = "%s"
    password = "%s"
  }
  target_cluster {
    endpoint = "%s"
    username = "%s"
    password = "%s"
  }
}
`, srcServer, srcUser, srcPass, drServer, drUser, drPass)
}
