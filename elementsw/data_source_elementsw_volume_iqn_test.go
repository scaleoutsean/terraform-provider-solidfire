package elementsw

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceElementSwVolumeIQN_basic(t *testing.T) {
	if os.Getenv("SOLIDFIRE_ACC") == "" {
		t.Skip("SOLIDFIRE_ACC must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceElementSwVolumeIQNConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.elementsw_volume_iqn.test", "iqn"),
					resource.TestCheckResourceAttrSet("data.elementsw_volume_iqn.test", "target_portal"),
				),
			},
		},
	})
}

const testAccDataSourceElementSwVolumeIQNConfig = `
data "elementsw_volume_iqn" "test" {
  unique_id  = "xh67"
  name       = "myvol"
  volume_id  = 1234
  svip       = "192.168.105.34"
}
`
