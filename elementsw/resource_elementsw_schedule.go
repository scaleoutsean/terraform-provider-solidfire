package elementsw

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func resourceElementswSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementswScheduleCreate,
		Read:   resourceElementswScheduleRead,
		Update: resourceElementswScheduleUpdate,
		Delete: resourceElementswScheduleDelete,
		Schema: map[string]*schema.Schema{
			"schedule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"hours": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"schedule_info": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"paused": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"recurring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"run_next_interval": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"starting_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"monthdays": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceElementswScheduleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	req := sdk.CreateScheduleRequest{
		ScheduleName:    d.Get("schedule_name").(string),
		ScheduleType:    d.Get("schedule_type").(string),
		Hours:           int64(d.Get("hours").(int)),
		Minutes:         int64(d.Get("minutes").(int)),
		Paused:          d.Get("paused").(bool),
		Recurring:       d.Get("recurring").(bool),
		RunNextInterval: d.Get("run_next_interval").(bool),
	}

	if v, ok := d.GetOk("starting_date"); ok {
		req.StartingDate = v.(string)
	}

	if v, ok := d.GetOk("attributes"); ok {
		req.Attributes = v
	}

	if v, ok := d.GetOk("schedule_info"); ok {
		m := v.(map[string]interface{})
		info := sdk.ScheduleInfo{}
		if val, ok := m["volumeID"]; ok {
			id, _ := strconv.ParseInt(val.(string), 10, 64)
			info.VolumeID = id
		}
		if val, ok := m["retention"]; ok {
			info.Retention = val.(string)
		}
		req.ScheduleInfo = info
	}

	if v, ok := d.GetOk("monthdays"); ok {
		days := v.([]interface{})
		intDays := make([]int64, len(days))
		for i, day := range days {
			intDays[i] = int64(day.(int))
		}
		req.Monthdays = intDays
	}

	id, err := client.CreateSchedule(&req)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", id))
	return resourceElementswScheduleRead(d, m)
}

func resourceElementswScheduleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	s, err := client.GetSchedule(id)
	if err != nil {
		return err
	}
	if s == nil {
		d.SetId("")
		return nil
	}

	d.Set("schedule_name", s.ScheduleName)
	d.Set("schedule_type", s.ScheduleType)
	d.Set("hours", int(s.Hours))
	d.Set("minutes", int(s.Minutes))
	d.Set("paused", s.Paused)
	d.Set("recurring", s.Recurring)
	d.Set("run_next_interval", s.RunNextInterval)
	d.Set("starting_date", s.StartingDate)
	d.Set("monthdays", s.Monthdays)
	return nil
}

func resourceElementswScheduleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	req := sdk.ModifyScheduleRequest{
		ScheduleID: id,
	}

	if d.HasChange("schedule_name") {
		req.ScheduleName = d.Get("schedule_name").(string)
	}
	if d.HasChange("paused") {
		req.Paused = d.Get("paused").(bool)
	}

	err := client.ModifySchedule(&req)
	if err != nil {
		return err
	}
	return resourceElementswScheduleRead(d, m)
}

func resourceElementswScheduleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	req := sdk.ModifyScheduleRequest{
		ScheduleID:  id,
		ToBeDeleted: true,
	}
	return client.ModifySchedule(&req)
}
