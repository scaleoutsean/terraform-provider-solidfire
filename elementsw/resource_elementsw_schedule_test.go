package elementsw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccElementswSchedule_basic(t *testing.T) {
	resourceName := "elementsw_schedule.test"
	scheduleName := "tfacc-schedule"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduleConfig(scheduleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schedule_name", scheduleName),
					resource.TestCheckResourceAttr(resourceName, "schedule_type", "snapshot"),
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

func testAccScheduleConfig(name string) string {
	return fmt.Sprintf(`
resource "elementsw_schedule" "test" {
  schedule_name = "%s"
  schedule_type = "snapshot"
  attributes = {
    frequency = "Time Interval"
  }
  minutes = 10
  schedule_info = {
    retention = "0:10:00"
    volumeID = 1
  }
  paused = false
  recurring = true
}
`, name)
}

// Additional tests for List, Get, Modify can be added here as needed.
