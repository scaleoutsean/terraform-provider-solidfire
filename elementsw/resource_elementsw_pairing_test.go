package elementsw

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccElementswVolumePairing_automated(t *testing.T) {
	timestamp := time.Now().Unix()
	timestampStr := fmt.Sprintf("%d", timestamp)
	drServer := os.Getenv("ELEMENTSW_SERVER_DR")
	if drServer == "" {
		t.Skip("ELEMENTSW_SERVER_DR not set, skipping replication tests")
	}

	srcServer := os.Getenv("ELEMENTSW_SERVER")
	if srcServer == "" {
		srcServer = "192.168.1.30"
	}
	if srcServer == drServer {
		t.Fatalf("srcServer and drServer are both %s", srcServer)
	}
	srcUser := os.Getenv("ELEMENTSW_USERNAME")
	srcPass := os.Getenv("ELEMENTSW_PASSWORD")
	srcVer := os.Getenv("ELEMENTSW_API_VERSION")
	if srcVer == "" {
		srcVer = "12.5"
	}

	drUser := os.Getenv("ELEMENTSW_USERNAME_DR")
	if drUser == "" {
		drUser = srcUser
	}
	drPass := os.Getenv("ELEMENTSW_PASSWORD_DR")
	if drPass == "" {
		drPass = srcPass
	}
	drVer := os.Getenv("ELEMENTSW_API_VERSION_DR")
	if drVer == "" {
		drVer = srcVer
	}

	log.Printf("[INFO] TEST PARAMS: src=%s, dr=%s", srcServer, drServer)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumePairingConfig(timestampStr, srcServer, srcVer, srcUser, srcPass, drServer, drVer, drUser, drPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("elementsw_volume_pairing.test", "pairing_key"),
					testAccVerifyVolumePairingBothSides("elementsw_volume.src", "elementsw_volume.dr"),
				),
			},
		},
	})
}

func testAccVerifyVolumePairingBothSides(srcRes, drRes string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		src, ok := s.RootModule().Resources[srcRes]
		if !ok {
			return fmt.Errorf("Source volume res not found: %s", srcRes)
		}
		dr, ok := s.RootModule().Resources[drRes]
		if !ok {
			return fmt.Errorf("DR volume res not found: %s", drRes)
		}

		srcID, _ := strconv.ParseInt(src.Primary.ID, 10, 64)
		drID, _ := strconv.ParseInt(dr.Primary.ID, 10, 64)

		// 1. Verify Source is readWrite and paired
		srcConfig := configStuct{
			User:            os.Getenv("ELEMENTSW_USERNAME"),
			Password:        os.Getenv("ELEMENTSW_PASSWORD"),
			ElementSwServer: os.Getenv("ELEMENTSW_SERVER"),
			APIVersion:      os.Getenv("ELEMENTSW_API_VERSION"),
		}
		if srcConfig.APIVersion == "" {
			srcConfig.APIVersion = "12.5"
		}
		srcClient, err := srcConfig.clientFun()
		if err != nil {
			return fmt.Errorf("failed to create src client: %w", err)
		}

		srcVol, err := srcClient.GetVolume(srcID)
		if err != nil {
			return fmt.Errorf("failed to get src vol: %w", err)
		}
		if srcVol.Access != "readWrite" {
			return fmt.Errorf("expected src volume %d to be readWrite, got %s", srcID, srcVol.Access)
		}

		// 2. Verify DR is replicationTarget
		drServer := os.Getenv("ELEMENTSW_SERVER_DR")
		drUser := os.Getenv("ELEMENTSW_USERNAME_DR")
		if drUser == "" {
			drUser = os.Getenv("ELEMENTSW_USERNAME")
		}
		drPass := os.Getenv("ELEMENTSW_PASSWORD_DR")
		if drPass == "" {
			drPass = os.Getenv("ELEMENTSW_PASSWORD")
		}
		drConfig := configStuct{
			User:            drUser,
			Password:        drPass,
			ElementSwServer: drServer,
			APIVersion:      os.Getenv("ELEMENTSW_API_VERSION_DR"),
		}
		if drConfig.APIVersion == "" {
			drConfig.APIVersion = srcConfig.APIVersion
		}
		drClient, err := drConfig.clientFun()
		if err != nil {
			return fmt.Errorf("failed to create dr client: %w", err)
		}

		drVol, err := drClient.GetVolume(drID)
		if err != nil {
			return fmt.Errorf("failed to get dr vol: %w", err)
		}
		if drVol.Access != "replicationTarget" {
			return fmt.Errorf("expected dr volume %d to be replicationTarget, got %s", drID, drVol.Access)
		}

		// 3. Verify pairing status on src
		vols, err := srcClient.ListActivePairedVolumes()
		if err != nil {
			return fmt.Errorf("failed to list paired volumes on src: %w", err)
		}

		found := false
		for _, v := range vols {
			if v.VolumeID == srcID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("volume %d not found in paired volumes on src", srcID)
		}

		log.Printf("[INFO] Volume pairing verified on both sides. Final wait for sync...")
		time.Sleep(10 * time.Second)

		return nil
	}
}

func testAccVolumePairingConfig(timestamp string, srcServer, srcVer, srcUser, srcPass, drServer, drVer, drUser, drPass string) string {
	drEndpoint := fmt.Sprintf("https://%s/json-rpc/%s", drServer, drVer)

	hcl := fmt.Sprintf(`
provider "elementsw" {
  elementsw_server = "%s"
  api_version      = "%s"
  username         = "%s"
  password         = "%s"
}

provider "elementswremote" {
  elementsw_server = "%s"
  api_version      = "%s"
  username         = "%s"
  password         = "%s"
}

resource "elementsw_account" "src" {
  username = "terraform-%s-src"
}

resource "elementsw_account" "dr" {
  provider = elementswremote
  username = "terraform-%s-dr"
}

resource "elementsw_volume" "src" {
  name       = "terraform-%s-vol"
  account_id = elementsw_account.src.account_id
  total_size = 10000000000
  enable512e = true
}

resource "elementsw_volume" "dr" {
  provider   = elementswremote
  name       = "terraform-%s-vol" 
  account_id = elementsw_account.dr.account_id
  total_size = 10000000000
  enable512e = true
}

resource "elementsw_volume_pairing" "test" {
  volume_id  = elementsw_volume.src.id
  mode       = "Async"
  
  target_cluster {
    endpoint = "%s"
    username = "%s"
    password = "%s"
  }

  depends_on = [elementsw_volume.dr]
}
`, srcServer, srcVer, srcUser, srcPass, drServer, drVer, drUser, drPass, timestamp, timestamp, timestamp, timestamp, drEndpoint, drUser, drPass)

	log.Printf("[INFO] Generated HCL:\n%s", hcl)
	return hcl
}
