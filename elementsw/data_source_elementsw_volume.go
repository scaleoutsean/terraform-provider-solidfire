package elementsw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func dataSourceElementSwVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementSwVolumeRead,
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"iqn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceElementSwVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	var foundVol *sdk.Volume

	if v, ok := d.GetOk("volume_id"); ok {
		vol, err := client.GetVolume(int64(v.(int)))
		if err != nil {
			return err
		}
		foundVol = vol
	} else if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		// List all active volumes and filter.
		// For a more efficient way, we'd use ListActiveVolumes with paging.
		client.initOnce.Do(client.init)
		startID := int64(0)
		for {
			resp, sdkErr := client.sdkClient.ListActiveVolumes(context.TODO(), &sdk.ListActiveVolumesRequest{
				StartVolumeID: startID,
			})
			if sdkErr != nil {
				return sdkErr
			}
			if len(resp.Volumes) == 0 {
				break
			}
			for _, vol := range resp.Volumes {
				if vol.Name == name {
					foundVol = &vol
					break
				}
				if vol.VolumeID > startID {
					startID = vol.VolumeID
				}
			}
			if foundVol != nil || len(resp.Volumes) < 1000 {
				break
			}
			startID++
		}
	} else {
		return fmt.Errorf("either volume_id or name must be specified")
	}

	if foundVol == nil {
		return fmt.Errorf("volume not found")
	}

	d.SetId(strconv.FormatInt(foundVol.VolumeID, 10))
	d.Set("volume_id", int(foundVol.VolumeID))
	d.Set("name", foundVol.Name)
	d.Set("account_id", int(foundVol.AccountID))
	d.Set("total_size", int(foundVol.TotalSize))
	d.Set("iqn", foundVol.Iqn)
	d.Set("access", foundVol.Access)

	return nil
}
