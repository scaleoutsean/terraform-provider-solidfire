package elementsw

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func resourceElementswQoSPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementswQoSPolicyCreate,
		Read:   resourceElementswQoSPolicyRead,
		Update: resourceElementswQoSPolicyUpdate,
		Delete: resourceElementswQoSPolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"qos_policy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},
			"qos": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"burst_iops": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"burst_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"max_iops": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"min_iops": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"curve": {
							Type:     schema.TypeMap,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Create QoS Policy
func resourceElementswQoSPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	name := d.Get("name").(string)
	qosList := d.Get("qos").([]interface{})
	if len(qosList) == 0 {
		return fmt.Errorf("qos block must be provided")
	}
	qosMap := qosList[0].(map[string]interface{})
	qos := sdk.QoS{
		MinIOPS:   int64(qosMap["min_iops"].(int)),
		MaxIOPS:   int64(qosMap["max_iops"].(int)),
		BurstIOPS: int64(qosMap["burst_iops"].(int)),
	}

	qosPolicyID, err := client.CreateQoSPolicy(name, qos)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", qosPolicyID))
	return resourceElementswQoSPolicyRead(d, meta)
}

// Update QoS Policy
func resourceElementswQoSPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	idStr := d.Id()
	qosPolicyID, _ := strconv.ParseInt(idStr, 10, 64)

	name := d.Get("name").(string)
	qosList := d.Get("qos").([]interface{})
	if len(qosList) == 0 {
		return fmt.Errorf("qos block must be provided")
	}
	qosMap := qosList[0].(map[string]interface{})
	qos := sdk.QoS{
		MinIOPS:   int64(qosMap["min_iops"].(int)),
		MaxIOPS:   int64(qosMap["max_iops"].(int)),
		BurstIOPS: int64(qosMap["burst_iops"].(int)),
	}

	err := client.ModifyQoSPolicy(qosPolicyID, name, qos)
	if err != nil {
		return err
	}
	return resourceElementswQoSPolicyRead(d, meta)
}

// Delete QoS Policy
func resourceElementswQoSPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	idStr := d.Id()
	qosPolicyID, _ := strconv.ParseInt(idStr, 10, 64)

	err := client.DeleteQoSPolicy(qosPolicyID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceElementswQoSPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	idStr := d.Id()
	qosPolicyID, _ := strconv.ParseInt(idStr, 10, 64)

	policy, err := client.GetQoSPolicy(qosPolicyID)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", policy.QosPolicyID))
	d.Set("name", policy.Name)
	d.Set("qos_policy_id", policy.QosPolicyID)
	d.Set("volume_ids", policy.VolumeIDs)

	qosMap := map[string]interface{}{
		"burst_iops": policy.Qos.BurstIOPS,
		"burst_time": policy.Qos.BurstTime,
		"max_iops":   policy.Qos.MaxIOPS,
		"min_iops":   policy.Qos.MinIOPS,
		"curve":      policy.Qos.Curve,
	}
	d.Set("qos", []interface{}{qosMap})

	return nil
}
