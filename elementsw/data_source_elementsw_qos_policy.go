package elementsw

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceElementSwQosPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementSwQosPolicyRead,
		Schema: map[string]*schema.Schema{
			"qos_policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"burst_iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceElementSwQosPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	var qosPolicyID int64
	var name string

	if v, ok := d.GetOk("qos_policy_id"); ok {
		qosPolicyID = int64(v.(int))
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if qosPolicyID == 0 && name == "" {
		return fmt.Errorf("one of qos_policy_id or name must be set")
	}

	policies, err := client.ListQoSPolicies()
	if err != nil {
		return fmt.Errorf("failed to list QoS policies: %w", err)
	}

	for _, p := range policies {
		if (qosPolicyID != 0 && p.QosPolicyID == qosPolicyID) || (name != "" && p.Name == name) {
			d.SetId(strconv.FormatInt(p.QosPolicyID, 10))
			d.Set("qos_policy_id", int(p.QosPolicyID))
			d.Set("name", p.Name)
			d.Set("min_iops", int(p.Qos.MinIOPS))
			d.Set("max_iops", int(p.Qos.MaxIOPS))
			d.Set("burst_iops", int(p.Qos.BurstIOPS))
			return nil
		}
	}

	return fmt.Errorf("QoS policy not found")
}
