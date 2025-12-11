package elementsw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccElementswSnapshot_basic(t *testing.T) {
	resourceName := "elementsw_snapshot.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-snap"),
					resource.TestCheckResourceAttrSet(resourceName, "created_snapshot_id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
		},
	})
}

func testAccSnapshotConfigBasic() string {
	return `
resource "elementsw_snapshot" "test" {
  volume_id = 1
  name = "test-snap"
}
`
}

func TestAccElementswGroupSnapshot_basic(t *testing.T) {
	resourceName := "elementsw_snapshot.group"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupSnapshotConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-group-snap"),
					resource.TestCheckResourceAttrSet(resourceName, "created_group_snapshot_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_group_snapshot_uuid"),
				),
			},
		},
	})
}

func testAccGroupSnapshotConfigBasic() string {
	return `
resource "elementsw_snapshot" "group" {
  volume_ids = [1, 2]
  name = "test-group-snap"
}
`
}
