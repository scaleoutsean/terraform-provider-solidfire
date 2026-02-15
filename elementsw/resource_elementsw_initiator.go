package elementsw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func resourceElementSwInitiator() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwInitiatorCreate,
		Read:   resourceElementSwInitiatorRead,
		Update: resourceElementSwInitiatorUpdate,
		Delete: resourceElementSwInitiatorDelete,
		Exists: resourceElementSwInitiatorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"volume_access_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iqns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceElementSwInitiatorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	req := sdk.CreateInitiatorsRequest{}
	newInit := sdk.CreateInitiator{}

	if v, ok := d.GetOk("name"); ok {
		newInit.Name = v.(string)
	}

	if v, ok := d.GetOk("alias"); ok {
		newInit.Alias = v.(string)
	}

	if v, ok := d.GetOk("volume_access_group_id"); ok {
		newInit.VolumeAccessGroupID = int64(v.(int))
	}

	// Attributes and other fields could be added here if supported by CreateInitiator struct
	req.Initiators = []sdk.CreateInitiator{newInit}

	client.initOnce.Do(client.init)
	res, sdkErr := client.sdkClient.CreateInitiators(context.TODO(), &req)
	if sdkErr != nil {
		return sdkErr
	}

	d.SetId(fmt.Sprintf("%v", res.Initiators[0].InitiatorID))
	return resourceElementSwInitiatorRead(d, meta)
}

func resourceElementSwInitiatorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	idStr := d.Id()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	req := sdk.ListInitiatorsRequest{
		Initiators: []int64{id},
	}
	client.initOnce.Do(client.init)
	res, sdkErr := client.sdkClient.ListInitiators(context.TODO(), &req)
	if sdkErr != nil {
		return sdkErr
	}

	if len(res.Initiators) != 1 {
		return fmt.Errorf("expected one initiator, got %d", len(res.Initiators))
	}

	init := res.Initiators[0]
	d.Set("name", init.InitiatorName)
	d.Set("alias", init.Alias)
	d.Set("attributes", init.Attributes)
	if len(init.VolumeAccessGroups) > 0 {
		d.Set("volume_access_group_id", init.VolumeAccessGroups[0])
	}

	return nil
}

func resourceElementSwInitiatorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	idStr := d.Id()
	id, _ := strconv.ParseInt(idStr, 10, 64)

	req := sdk.ModifyInitiatorsRequest{}
	modInit := sdk.ModifyInitiator{
		InitiatorID: id,
	}

	if d.HasChange("alias") {
		modInit.Alias = d.Get("alias").(string)
	}
	if d.HasChange("volume_access_group_id") {
		modInit.VolumeAccessGroupID = int64(d.Get("volume_access_group_id").(int))
	}

	req.Initiators = []sdk.ModifyInitiator{modInit}

	client.initOnce.Do(client.init)
	_, sdkErr := client.sdkClient.ModifyInitiators(context.TODO(), &req)
	return sdkErr
}

func resourceElementSwInitiatorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	idStr := d.Id()
	id, _ := strconv.ParseInt(idStr, 10, 64)

	req := sdk.DeleteInitiatorsRequest{
		Initiators: []int64{id},
	}
	client.initOnce.Do(client.init)
	_, sdkErr := client.sdkClient.DeleteInitiators(context.TODO(), &req)
	return sdkErr
}

func resourceElementSwInitiatorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client)
	idStr := d.Id()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return false, nil
	}

	req := sdk.ListInitiatorsRequest{
		Initiators: []int64{id},
	}
	client.initOnce.Do(client.init)
	res, sdkErr := client.sdkClient.ListInitiators(context.TODO(), &req)
	if sdkErr != nil {
		// Check for 500:xUnknown or similar
		if sdkErr.Detail == "500:xUnknown" {
			return false, nil
		}
		return false, sdkErr
	}

	return len(res.Initiators) == 1, nil
}
