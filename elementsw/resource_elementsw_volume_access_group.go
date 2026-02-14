package elementsw

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func resourceElementSwVolumeAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwVolumeAccessGroupCreate,
		Read:   resourceElementSwVolumeAccessGroupRead,
		Update: resourceElementSwVolumeAccessGroupUpdate,
		Delete: resourceElementSwVolumeAccessGroupDelete,
		Exists: resourceElementSwVolumeAccessGroupExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"initiators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceElementSwVolumeAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume access group: %#v", d)
	client := meta.(*Client)

	req := sdk.CreateVolumeAccessGroupRequest{}

	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
	}

	if raw, ok := d.GetOk("volumes"); ok {
		for _, v := range raw.([]interface{}) {
			req.Volumes = append(req.Volumes, int64(v.(int)))
		}
	}

	client.initOnce.Do(client.init)
	res, sdkErr := client.sdkClient.CreateVolumeAccessGroup(context.TODO(), &req)
	if sdkErr != nil {
		return sdkErr
	}

	d.SetId(fmt.Sprintf("%v", res.VolumeAccessGroupID))
	return resourceElementSwVolumeAccessGroupRead(d, meta)
}

func resourceElementSwVolumeAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume access group: %#v", d)
	client := meta.(*Client)

	idStr := d.Id()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	req := sdk.ListVolumeAccessGroupsRequest{
		VolumeAccessGroups: []int64{id},
	}
	client.initOnce.Do(client.init)
	res, sdkErr := client.sdkClient.ListVolumeAccessGroups(context.TODO(), &req)
	if sdkErr != nil {
		return sdkErr
	}

	if len(res.VolumeAccessGroups) != 1 {
		return fmt.Errorf("expected one volume access group")
	}

	vag := res.VolumeAccessGroups[0]
	d.Set("name", vag.Name)
	d.Set("initiators", vag.Initiators)
	d.Set("volumes", vag.Volumes)

	return nil
}

func resourceElementSwVolumeAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	idStr := d.Id()
	id, _ := strconv.ParseInt(idStr, 10, 64)

	req := sdk.ModifyVolumeAccessGroupRequest{
		VolumeAccessGroupID: id,
	}

	if d.HasChange("name") {
		req.Name = d.Get("name").(string)
	}

	if d.HasChange("volumes") {
		raw := d.Get("volumes").([]interface{})
		for _, v := range raw {
			req.Volumes = append(req.Volumes, int64(v.(int)))
		}
	}

	client.initOnce.Do(client.init)
	_, sdkErr := client.sdkClient.ModifyVolumeAccessGroup(context.TODO(), &req)
	return sdkErr
}

func resourceElementSwVolumeAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	idStr := d.Id()
	id, _ := strconv.ParseInt(idStr, 10, 64)

	req := sdk.DeleteVolumeAccessGroupRequest{
		VolumeAccessGroupID: id,
	}
	client.initOnce.Do(client.init)
	_, sdkErr := client.sdkClient.DeleteVolumeAccessGroup(context.TODO(), &req)
	return sdkErr
}

func resourceElementSwVolumeAccessGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client)
	idStr := d.Id()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return false, nil
	}

	req := sdk.ListVolumeAccessGroupsRequest{
		VolumeAccessGroups: []int64{id},
	}
	client.initOnce.Do(client.init)
	res, sdkErr := client.sdkClient.ListVolumeAccessGroups(context.TODO(), &req)
	if sdkErr != nil {
		if sdkErr.Detail == "500:xUnknown" {
			return false, nil
		}
		return false, sdkErr
	}

	return len(res.VolumeAccessGroups) == 1, nil
}
