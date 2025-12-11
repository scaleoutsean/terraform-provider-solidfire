package elementsw

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	qos := qosDetails{
		MinIOPS:   qosMap["min_iops"].(int),
		MaxIOPS:   qosMap["max_iops"].(int),
		BurstIOPS: qosMap["burst_iops"].(int),
		BurstTime: qosMap["burst_time"].(int),
		Curve:     map[string]int{}, // Curve can be set if needed
	}
	req := createQoSPolicyRequest{Name: name, QoS: qos}
	result, err := client.CreateQoSPolicy(req)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", result.QoSPolicyID))
	d.Set("qos_policy_id", result.QoSPolicyID)
	return resourceElementswQoSPolicyRead(d, meta)
}

// Update QoS Policy
func resourceElementswQoSPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	qosPolicyID := d.Get("qos_policy_id").(int)
	name := ""
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	qosList := d.Get("qos").([]interface{})
	if len(qosList) == 0 {
		return fmt.Errorf("qos block must be provided")
	}
	qosMap := qosList[0].(map[string]interface{})
	qos := qosDetails{
		MinIOPS:   qosMap["min_iops"].(int),
		MaxIOPS:   qosMap["max_iops"].(int),
		BurstIOPS: qosMap["burst_iops"].(int),
		BurstTime: qosMap["burst_time"].(int),
		Curve:     map[string]int{}, // Curve can be set if needed
	}
	req := modifyQoSPolicyRequest{QoSPolicyID: qosPolicyID, Name: name, QoS: qos}
	_, err := client.ModifyQoSPolicy(req)
	if err != nil {
		return err
	}
	return resourceElementswQoSPolicyRead(d, meta)
}

// Delete QoS Policy
func resourceElementswQoSPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	qosPolicyID := d.Get("qos_policy_id").(int)
	req := deleteQoSPolicyRequest{QoSPolicyID: qosPolicyID}
	_, err := client.DeleteQoSPolicy(req)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceElementswQoSPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	qosPolicyID, ok := d.Get("qos_policy_id").(int)
	if !ok || qosPolicyID == 0 {
		return fmt.Errorf("qos_policy_id must be set")
	}

	result, err := client.GetQoSPolicy(getQoSPolicyRequest{QoSPolicyID: qosPolicyID})
	if err != nil {
		return err
	}

	policy := result.QoSPolicy
	d.SetId(fmt.Sprintf("%d", policy.QoSPolicyID))
	d.Set("name", policy.Name)
	d.Set("qos_policy_id", policy.QoSPolicyID)
	d.Set("volume_ids", policy.VolumeIDs)

	qosMap := map[string]interface{}{
		"burst_iops": policy.QoS.BurstIOPS,
		"burst_time": policy.QoS.BurstTime,
		"max_iops":   policy.QoS.MaxIOPS,
		"min_iops":   policy.QoS.MinIOPS,
		"curve":      policy.QoS.Curve,
	}
	d.Set("qos", []interface{}{qosMap})

	return nil
}
