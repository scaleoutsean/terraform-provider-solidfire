package elementsw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func dataSourceElementSwVolumeAccessGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementSwVolumeAccessGroupRead,
		Schema: map[string]*schema.Schema{
			"volume_access_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"initiators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"attributes": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceElementSwVolumeAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	var vagID int64
	name, hasName := d.GetOk("name")

	if v, ok := d.GetOk("volume_access_group_id"); ok {
		vagID = int64(v.(int))
	}

	if vagID == 0 && !hasName {
		return fmt.Errorf("one of volume_access_group_id or name must be set")
	}

	req := sdk.ListVolumeAccessGroupsRequest{}
	if vagID != 0 {
		req.VolumeAccessGroups = []int64{vagID}
	}

	client.initOnce.Do(client.init)
	res, sdkErr := client.sdkClient.ListVolumeAccessGroups(context.TODO(), &req)
	if sdkErr != nil {
		return fmt.Errorf("failed to list volume access groups: %s", sdkErr.Detail)
	}

	for _, vag := range res.VolumeAccessGroups {
		if (vagID != 0 && vag.VolumeAccessGroupID == vagID) || (hasName && vag.Name == name.(string)) {
			d.SetId(strconv.FormatInt(vag.VolumeAccessGroupID, 10))
			d.Set("volume_access_group_id", int(vag.VolumeAccessGroupID))
			d.Set("name", vag.Name)
			d.Set("initiators", vag.Initiators)

			vols := make([]int, len(vag.Volumes))
			for i, v := range vag.Volumes {
				vols[i] = int(v)
			}
			d.Set("volumes", vols)
			// d.Set("attributes", vag.Attributes)
			return nil
		}
	}

	return fmt.Errorf("volume access group not found")
}
