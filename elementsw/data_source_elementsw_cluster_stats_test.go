package elementsw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceElementSwClusterStats_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceElementSwClusterStatsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.elementsw_cluster_stats.test", "volume_count"),
					resource.TestCheckResourceAttrSet("data.elementsw_cluster_stats.test", "node_count"),
					resource.TestCheckResourceAttrSet("data.elementsw_cluster_stats.test", "capacity.0.used_space"),
					resource.TestCheckResourceAttrSet("data.elementsw_cluster_stats.test", "metrics.0.actual_iops"),
				),
			},
		},
	})
}

const testAccDataSourceElementSwClusterStatsConfig = `
data "elementsw_cluster_stats" "test" {}
`
