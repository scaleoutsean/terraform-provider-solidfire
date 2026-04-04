package solidfire

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccElementswSnapshot_basic(t *testing.T) {
	resourceName := "solidfire_snapshot.test"
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
resource "solidfire_account" "test" {
  username = "tf-acc-test-snap"
}

resource "solidfire_volume" "test" {
  name = "tf-acc-test-snap-vol"
  account_id = solidfire_account.test.id
  total_size = 1073741824
  enable512e = true
}

resource "solidfire_snapshot" "test" {
  volume_id = solidfire_volume.test.id
  name = "test-snap"
}
`
}

func TestAccElementswGroupSnapshot_basic(t *testing.T) {
	resourceName := "solidfire_snapshot.group"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupSnapshotConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-group-snap"),
					resource.TestCheckResourceAttrSet(resourceName, "created_group_snapshot_id"),
				),
			},
		},
	})
}

func testAccGroupSnapshotConfigBasic() string {
	return `
resource "solidfire_account" "test" {
  username = "tf-acc-test-snap"
}

resource "solidfire_volume" "test" {
  name = "tf-acc-test-grp-snap-vol"
  account_id = solidfire_account.test.id
  total_size = 1073741824
  enable512e = true
}

resource "solidfire_snapshot" "group" {
  volume_ids = [solidfire_volume.test.id]
  name = "test-group-snap"
}
`
}
