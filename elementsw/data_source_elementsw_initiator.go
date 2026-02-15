package elementsw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func dataSourceElementSwInitiator() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementSwInitiatorRead,
		Schema: map[string]*schema.Schema{
			"initiator_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_access_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceElementSwInitiatorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	var initiatorID int64
	var name string

	if v, ok := d.GetOk("initiator_id"); ok {
		initiatorID = int64(v.(int))
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if initiatorID == 0 && name == "" {
		return fmt.Errorf("one of initiator_id or name must be set")
	}

	req := sdk.ListInitiatorsRequest{}
	if initiatorID != 0 {
		req.Initiators = []int64{initiatorID}
	}

	client.initOnce.Do(client.init)
	res, sdkErr := client.sdkClient.ListInitiators(context.TODO(), &req)
	if sdkErr != nil {
		return fmt.Errorf("failed to list initiators: %s", sdkErr.Detail)
	}

	for _, init := range res.Initiators {
		if (initiatorID != 0 && init.InitiatorID == initiatorID) || (name != "" && init.InitiatorName == name) {
			d.SetId(strconv.FormatInt(init.InitiatorID, 10))
			d.Set("initiator_id", int(init.InitiatorID))
			d.Set("name", init.InitiatorName)
			d.Set("alias", init.Alias)
			vags := make([]int, len(init.VolumeAccessGroups))
			for i, vagID := range init.VolumeAccessGroups {
				vags[i] = int(vagID)
			}
			d.Set("volume_access_group_ids", vags)
			return nil
		}
	}

	return fmt.Errorf("initiator not found")
}
