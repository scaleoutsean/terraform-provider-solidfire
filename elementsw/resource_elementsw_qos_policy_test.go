package elementsw

import (
	"testing"
	// "fmt" // Removed unused import
	// "os"  // Removed unused import

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccElementswQoSPolicy_CRUD(t *testing.T) {
	resourceName := "elementsw_qos_policy.example"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { /* add pre-checks if needed */ },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "elementsw_qos_policy" "example" {
									   name = "test-qos-policy"
									   qos {
											   min_iops   = 100
											   max_iops   = 200
											   burst_iops = 300
											   burst_time = 60
									   }
							   }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "qos_policy_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test-qos-policy"),
					resource.TestCheckResourceAttr(resourceName, "qos.0.burst_iops", "300"),
				),
			},
			{
				Config: `resource "elementsw_qos_policy" "example" {
									   name = "updated-qos-policy"
									   qos {
											   min_iops   = 150
											   max_iops   = 250
											   burst_iops = 350
											   burst_time = 120
									   }
							   }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "updated-qos-policy"),
					resource.TestCheckResourceAttr(resourceName, "qos.0.min_iops", "150"),
					resource.TestCheckResourceAttr(resourceName, "qos.0.burst_time", "120"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
