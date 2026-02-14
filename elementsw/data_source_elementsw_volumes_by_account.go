package elementsw

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceElementswVolumesByAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementswVolumesByAccountRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"volume_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceElementswVolumesByAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	accountID := int64(d.Get("account_id").(int))

	volumes, err := client.ListVolumesForAccount(accountID)
	if err != nil {
		return fmt.Errorf("failed to list volumes for account: %w", err)
	}
	var ids []int64
	for _, v := range volumes {
		ids = append(ids, v.VolumeID)
	}
	d.SetId(fmt.Sprintf("%d", accountID))
	d.Set("volume_ids", ids)
	return nil
}
