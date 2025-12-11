package elementsw

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceElementSwCluster_basic(t *testing.T) {
	if os.Getenv("SOLIDFIRE_ACC") == "" {
		t.Skip("SOLIDFIRE_ACC must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceElementSwClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.elementsw_cluster.test", "name"),
					resource.TestCheckResourceAttrSet("data.elementsw_cluster.test", "unique_id"),
					resource.TestCheckResourceAttrSet("data.elementsw_cluster.test", "cluster_version"),
					resource.TestCheckResourceAttrSet("data.elementsw_cluster.test", "cluster_api_version"),
					resource.TestCheckResourceAttrSet("data.elementsw_cluster.test", "mvip"),
					resource.TestCheckResourceAttrSet("data.elementsw_cluster.test", "svip"),
				),
			},
		},
	})
}

const testAccDataSourceElementSwClusterConfig = `
data "elementsw_cluster" "test" {}
`
