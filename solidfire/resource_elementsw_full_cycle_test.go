package solidfire

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccElementsw_FullCycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create Account and QoS Policy
			{
				Config: testAccFullCycleConfig_Step1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidfire_account.test", "id"),
					resource.TestCheckResourceAttr("solidfire_account.test", "username", "tf-acc-test-account-2"),
					resource.TestCheckResourceAttrSet("solidfire_qos_policy.test", "id"),
					resource.TestCheckResourceAttr("solidfire_qos_policy.test", "name", "tf-acc-test-policy"),
				),
			},
			// Step 2: Add Volume
			{
				Config: testAccFullCycleConfig_Step2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidfire_volume.test", "id"),
					resource.TestCheckResourceAttr("solidfire_volume.test", "name", "tf-acc-test-volume"),
					resource.TestCheckResourceAttr("solidfire_volume.test", "total_size", "1073741824"),
					testAccCheckVolumeQoSPolicyID("solidfire_volume.test", "solidfire_qos_policy.test"),
				),
			},
			// Step 3: Add Volume Access Group
			{
				Config: testAccFullCycleConfig_Step3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidfire_volume_access_group.test", "id"),
					resource.TestCheckResourceAttr("solidfire_volume_access_group.test", "name", "tf-acc-test-vag"),
					resource.TestCheckResourceAttr("solidfire_volume_access_group.test", "volumes.#", "1"),
					resource.TestCheckResourceAttrSet("solidfire_initiator.test", "id"),
				),
			},
			// Step 4: Add Schedule
			{
				Config: testAccFullCycleConfig_Step4,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidfire_schedule.test", "id"),
					resource.TestCheckResourceAttr("solidfire_schedule.test", "schedule_name", "tf-acc-test-schedule"),
				),
			},
			// Step 5: Update Volume Size
			{
				Config: testAccFullCycleConfig_Step5,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("solidfire_volume.test", "total_size", "2000683008"),
				),
			},
			// Step 6: Update QoS Policy
			{
				Config: testAccFullCycleConfig_Step6,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("solidfire_qos_policy.test", "qos.0.min_iops", "100"),
				),
			},
		},
	})
}

func testAccCheckVolumeQoSPolicyID(volResource, qosResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vol, ok := s.RootModule().Resources[volResource]
		if !ok {
			return fmt.Errorf("Not found: %s", volResource)
		}
		qos, ok := s.RootModule().Resources[qosResource]
		if !ok {
			return fmt.Errorf("Not found: %s", qosResource)
		}

		if vol.Primary.Attributes["qos_policy_id"] != qos.Primary.ID {
			return fmt.Errorf("Volume QoS Policy ID %s does not match QoS Policy ID %s",
				vol.Primary.Attributes["qos_policy_id"], qos.Primary.ID)
		}
		return nil
	}
}

const testAccFullCycleConfig_Step1 = `
resource "solidfire_account" "test" {
  username = "tf-acc-test-account-2"
}

resource "solidfire_qos_policy" "test" {
  name = "tf-acc-test-policy"
  qos {
    min_iops = 50
    max_iops = 1000
    burst_iops = 2000
  }
}
`

const testAccFullCycleConfig_Step2 = testAccFullCycleConfig_Step1 + `
resource "solidfire_volume" "test" {
  name = "tf-acc-test-volume"
  account_id = solidfire_account.test.id
  total_size = 1073741824
  enable512e = true
  qos_policy_id = solidfire_qos_policy.test.id
}
`

const testAccFullCycleConfig_Step3 = testAccFullCycleConfig_Step2 + `
resource "solidfire_volume_access_group" "test" {
  name = "tf-acc-test-vag"
  volumes = [solidfire_volume.test.id]
}

resource "solidfire_initiator" "test" {
  name = "tf-acc-test-initiator"
  volume_access_group_id = solidfire_volume_access_group.test.id
}
`

const testAccFullCycleConfig_Step4 = testAccFullCycleConfig_Step3 + `
resource "solidfire_schedule" "test" {
  schedule_name = "tf-acc-test-schedule"
  schedule_type = "Snapshot"
  attributes = {
    frequency = "Time Interval"
  }
  minutes = 60
  schedule_info = {
    volumeID = solidfire_volume.test.id
  }
}
`

const testAccFullCycleConfig_Step5 = testAccFullCycleConfig_Step1 + `
resource "solidfire_volume" "test" {
  name = "tf-acc-test-volume"
  account_id = solidfire_account.test.id
  total_size = 2000683008 # Increased size
  enable512e = true
  qos_policy_id = solidfire_qos_policy.test.id
}

resource "solidfire_volume_access_group" "test" {
  name = "tf-acc-test-vag"
  volumes = [solidfire_volume.test.id]
}

resource "solidfire_initiator" "test" {
  name = "tf-acc-test-initiator"
  volume_access_group_id = solidfire_volume_access_group.test.id
}

resource "solidfire_schedule" "test" {
  schedule_name = "tf-acc-test-schedule"
  schedule_type = "Snapshot"
  attributes = {
    frequency = "Time Interval"
  }
  minutes = 60
  schedule_info = {
    volumeID = solidfire_volume.test.id
  }
}
`

const testAccFullCycleConfig_Step6 = `
resource "solidfire_account" "test" {
  username = "tf-acc-test-account-2"
}

resource "solidfire_qos_policy" "test" {
  name = "tf-acc-test-policy"
  qos {
    min_iops = 100 # Increased Min IOPS
    max_iops = 1000
    burst_iops = 2000
  }
}

resource "solidfire_volume" "test" {
  name = "tf-acc-test-volume"
  account_id = solidfire_account.test.id
  total_size = 2000683008
  enable512e = true
  qos_policy_id = solidfire_qos_policy.test.id
}

resource "solidfire_volume_access_group" "test" {
  name = "tf-acc-test-vag"
  volumes = [solidfire_volume.test.id]
}

resource "solidfire_initiator" "test" {
  name = "tf-acc-test-initiator"
  volume_access_group_id = solidfire_volume_access_group.test.id
}

resource "solidfire_schedule" "test" {
  schedule_name = "tf-acc-test-schedule"
  schedule_type = "Snapshot"
  attributes = {
    frequency = "Time Interval"
  }
  minutes = 60
  schedule_info = {
    volumeID = solidfire_volume.test.id
  }
}
`
