package elementsw

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceElementSwVolumeIQN() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementSwVolumeIQNRead,
		Schema: map[string]*schema.Schema{
			"unique_id": {Type: schema.TypeString, Required: true},
			"name": {Type: schema.TypeString, Required: true},
			"volume_id": {Type: schema.TypeInt, Required: true},
			"svip": {Type: schema.TypeString, Required: true},
			"iqn": {Type: schema.TypeString, Computed: true},
			"target_portal": {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceElementSwVolumeIQNRead(d *schema.ResourceData, meta interface{}) error {
	uniqueID := d.Get("unique_id").(string)
	name := d.Get("name").(string)
	volumeID := d.Get("volume_id").(int)
	svip := d.Get("svip").(string)

	// Compose IQN
	// iqn.2010-01.com.solidfire:<uniqueID>.<name>.<volumeID>
	volumeIQN := fmt.Sprintf("iqn.2010-01.com.solidfire:%s.%s.%d", uniqueID, name, volumeID)

	d.SetId(volumeIQN)
	d.Set("iqn", volumeIQN)
	d.Set("target_portal", svip)

	return nil
}
