package elementsw

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func resourceElementSwVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwVolumeCreate,
		Read:   resourceElementSwVolumeRead,
		Update: resourceElementSwVolumeUpdate,
		Delete: resourceElementSwVolumeDelete,
		Exists: resourceElementSwVolumeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true, // we might resolve name to ID
			},
			"total_size": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1073741824 {
						errors = append(errors, fmt.Errorf("%q must be at least 1073741824 bytes (1 GiB)", k))
					}
					return
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					o, _ := strconv.ParseInt(old, 10, 64)
					n, _ := strconv.ParseInt(new, 10, 64)
					diff := o - n
					if diff < 0 {
						diff = -diff
					}
					return diff < 1048576 // Ignore differences less than 1MiB due to API rounding
				},
			},
			"enable512e": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"min_iops": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_iops": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"burst_iops": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"qos_policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"access": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					valid := map[string]bool{
						"readWrite":         true,
						"readOnly":          true,
						"locked":            true,
						"replicationTarget": true,
					}
					if !valid[value] {
						errors = append(errors, fmt.Errorf("%q is not a valid volume access mode", value))
					}
					return
				},
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"iqn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceElementSwVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	// Resolve account name to ID if needed
	accountID := int64(0)
	if v, ok := d.GetOk("account"); ok {
		acc, err := client.GetAccountByName(v.(string))
		if err != nil {
			return fmt.Errorf("failed to find account %s: %w", v.(string), err)
		}
		accountID = acc.AccountID
	} else if v, ok := d.GetOk("account_id"); ok {
		accountID = int64(v.(int))
	} else {
		return fmt.Errorf("either account or account_id must be provided")
	}

	req := sdk.CreateVolumeRequest{
		Name:       d.Get("name").(string),
		AccountID:  accountID,
		TotalSize:  int64(d.Get("total_size").(int)),
		Enable512e: d.Get("enable512e").(bool),
	}

	if v, ok := d.GetOk("access"); ok {
		req.Access = v.(string)
	}

	if v, ok := d.GetOk("qos_policy_id"); ok {
		req.QosPolicyID = int64(v.(int))
	} else {
		req.Qos = sdk.QoS{
			MinIOPS:   int64(d.Get("min_iops").(int)),
			MaxIOPS:   int64(d.Get("max_iops").(int)),
			BurstIOPS: int64(d.Get("burst_iops").(int)),
		}
	}

	// Attributes
	if v, ok := d.GetOk("attributes"); ok {
		req.Attributes = v.(map[string]interface{})
	}

	client.initOnce.Do(client.init)
	resp, sdkErr := client.sdkClient.CreateVolume(context.TODO(), &req)
	if sdkErr != nil {
		return fmt.Errorf("CreateVolume failed: %s", sdkErr.Detail)
	}

	d.SetId(fmt.Sprintf("%d", resp.VolumeID))
	return resourceElementSwVolumeRead(d, meta)
}

func resourceElementSwVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)

	vol, err := client.GetVolume(id)
	if err != nil {
		return err
	}

	d.Set("name", vol.Name)
	d.Set("account_id", int(vol.AccountID))
	d.Set("total_size", int(vol.TotalSize))
	d.Set("enable512e", vol.Enable512e)
	d.Set("iqn", vol.Iqn)
	d.Set("access", vol.Access)
	if vol.QosPolicyID != 0 {
		d.Set("qos_policy_id", int(vol.QosPolicyID))
	} else {
		d.Set("min_iops", int(vol.Qos.MinIOPS))
		d.Set("max_iops", int(vol.Qos.MaxIOPS))
		d.Set("burst_iops", int(vol.Qos.BurstIOPS))
	}
	// Attributes...
	return nil
}

func resourceElementSwVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)

	req := sdk.ModifyVolumeRequest{
		VolumeID: id,
	}

	if d.HasChange("total_size") {
		req.TotalSize = int64(d.Get("total_size").(int))
	}
	if d.HasChange("access") {
		req.Access = d.Get("access").(string)
	}
	if d.HasChange("qos_policy_id") {
		req.QosPolicyID = int64(d.Get("qos_policy_id").(int))
	} else if d.HasChange("min_iops") || d.HasChange("max_iops") || d.HasChange("burst_iops") {
		req.Qos = sdk.QoS{
			MinIOPS:   int64(d.Get("min_iops").(int)),
			MaxIOPS:   int64(d.Get("max_iops").(int)),
			BurstIOPS: int64(d.Get("burst_iops").(int)),
		}
	}

	client.initOnce.Do(client.init)
	_, sdkErr := client.sdkClient.ModifyVolume(context.TODO(), &req)
	if sdkErr != nil {
		return fmt.Errorf("ModifyVolume failed: %s", sdkErr.Detail)
	}

	return resourceElementSwVolumeRead(d, meta)
}

func resourceElementSwVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)

	client.initOnce.Do(client.init)
	_, sdkErr := client.sdkClient.DeleteVolume(context.TODO(), &sdk.DeleteVolumeRequest{VolumeID: id})
	if sdkErr != nil {
		return fmt.Errorf("DeleteVolume failed: %s", sdkErr.Detail)
	}

	_, sdkErr = client.sdkClient.PurgeDeletedVolume(context.TODO(), &sdk.PurgeDeletedVolumeRequest{VolumeID: id})
	if sdkErr != nil {
		log.Printf("[WARN] PurgeDeletedVolume failed for %d: %s", id, sdkErr.Detail)
	}

	d.SetId("")
	return nil
}

func resourceElementSwVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, nil
	}

	_, err = client.GetVolume(id)
	if err != nil {
		return false, nil
	}

	return true, nil
}
