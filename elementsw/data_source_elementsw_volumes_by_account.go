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
   accountID := d.Get("account_id").(int)
   req := listVolumesByAccountIDRequest{Accounts: []int{accountID}}
   result, err := client.listVolumesByVolumeID(req)
   if err != nil {
	   return fmt.Errorf("failed to list volumes for account: %w", err)
   }
   var ids []int
   for _, v := range result.Volumes {
	   ids = append(ids, v.VolumeID)
   }
   d.SetId(fmt.Sprintf("%d", accountID))
   d.Set("volume_ids", ids)
   return nil
}
