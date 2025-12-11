package elementsw

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// --- helpers for single vs group snapshot schedule ---
// Converts a slice of volume IDs to the correct scheduleInfo field for API requests

// Extracts a slice of volume IDs from scheduleInfo, regardless of single/group
func extractVolumeIDs(scheduleInfo map[string]interface{}) []string {
   if scheduleInfo == nil {
       return nil
   }
   if v, ok := scheduleInfo["volumes"]; ok {
       if arr, ok := v.([]interface{}); ok {
           ids := make([]string, len(arr))
           for i, id := range arr {
               ids[i] = id.(string)
           }
           return ids
       }
   }
   if v, ok := scheduleInfo["volumeID"]; ok {
       if id, ok := v.(string); ok {
           return []string{id}
       }
   }
   return nil
}

// resourceElementswSchedule returns the Terraform resource for SolidFire snapshot schedules
func resourceElementswSchedule() *schema.Resource {
   return &schema.Resource{
      Create: resourceElementswScheduleCreate,
      Read:   resourceElementswScheduleRead,
      Update: resourceElementswScheduleUpdate,
      Delete: resourceElementswScheduleDelete,
      Schema: map[string]*schema.Schema{
         "schedule_name":     {Type: schema.TypeString, Required: true},
         "schedule_type":     {Type: schema.TypeString, Required: true},
         "attributes":        {Type: schema.TypeMap, Optional: true},
         "minutes":           {Type: schema.TypeInt, Optional: true},
         "hours":             {Type: schema.TypeInt, Optional: true},
         "schedule_info":     {Type: schema.TypeMap, Optional: true},
         "paused":            {Type: schema.TypeBool, Optional: true},
         "recurring":         {Type: schema.TypeBool, Optional: true},
         "run_next_interval": {Type: schema.TypeBool, Optional: true},
         "starting_date":     {Type: schema.TypeString, Optional: true},
         "monthdays":         {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeInt}},
         "weekdays":          {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeMap}},
      },
   }
}

// resourceElementswScheduleCreate handles the Create operation
func resourceElementswScheduleCreate(d *schema.ResourceData, m interface{}) error {
   client := m.(*Client)
   // Use buildScheduleInfo if volumes are provided in attributes
   var scheduleInfo map[string]interface{}
   if v, ok := d.GetOk("schedule_info"); ok {
      scheduleInfo = toStringMap(v)
   }
   // Example: if you want to build scheduleInfo from volumes, retention, snapMirrorLabel
   // Uncomment and adapt as needed:
   // volumes := []string{"vol1", "vol2"} // get from d.Get("volumes") or elsewhere
   // retention := "" // get from d.Get("retention") or elsewhere
   // snapMirrorLabel := nil // get from d.Get("snapMirrorLabel") or elsewhere
   // scheduleInfo = buildScheduleInfo(volumes, retention, snapMirrorLabel)
   req := CreateScheduleRequest{
      ScheduleName: d.Get("schedule_name").(string),
      ScheduleType: d.Get("schedule_type").(string),
      Attributes:   toStringMap(d.Get("attributes")),
      ScheduleInfo: scheduleInfo,
   }
   if v, ok := d.GetOk("minutes"); ok { val := v.(int); req.Minutes = &val }
   if v, ok := d.GetOk("hours"); ok { val := v.(int); req.Hours = &val }
   if v, ok := d.GetOk("paused"); ok { val := v.(bool); req.Paused = &val }
   if v, ok := d.GetOk("recurring"); ok { val := v.(bool); req.Recurring = &val }
   if v, ok := d.GetOk("run_next_interval"); ok { val := v.(bool); req.RunNextInterval = &val }
   if v, ok := d.GetOk("starting_date"); ok { val := v.(string); req.StartingDate = &val }
   if v, ok := d.GetOk("monthdays"); ok { req.Monthdays = toIntSlice(v) }
   if v, ok := d.GetOk("weekdays"); ok { req.Weekdays = toWeekdaysSlice(v) }
   res, err := client.CreateSchedule(req)
   if err != nil {
      return err
   }
   d.SetId(fmt.Sprintf("%d", res.Schedule.ScheduleID))
   return resourceElementswScheduleRead(d, m)
}

// resourceElementswScheduleRead handles the Read operation
func resourceElementswScheduleRead(d *schema.ResourceData, m interface{}) error {
   client := m.(*Client)
   id, err := toIntID(d.Id())
   if err != nil {
      return err
   }
   res, err := client.GetSchedule(GetScheduleRequest{ScheduleID: id})
   if err != nil {
      return err
   }
   s := res.Schedule
   d.Set("schedule_name", s.ScheduleName)
   d.Set("schedule_type", s.ScheduleType)
   d.Set("attributes", s.Attributes)
   d.Set("minutes", s.Minutes)
   d.Set("hours", s.Hours)
   d.Set("paused", s.Paused)
   d.Set("recurring", s.Recurring)
   d.Set("run_next_interval", s.RunNextInterval)
   d.Set("schedule_info", s.ScheduleInfo)
   // Use extractVolumeIDs to demonstrate usage
   _ = extractVolumeIDs(s.ScheduleInfo)
   d.Set("starting_date", s.StartingDate)
   d.Set("monthdays", s.Monthdays)
   d.Set("weekdays", s.Weekdays)
   return nil
}

// resourceElementswScheduleUpdate handles the Update operation
func resourceElementswScheduleUpdate(d *schema.ResourceData, m interface{}) error {
   client := m.(*Client)
   id, err := toIntID(d.Id())
   if err != nil {
      return err
   }
   req := ModifyScheduleRequest{ScheduleID: id}
   if v, ok := d.GetOk("attributes"); ok { req.Attributes = toStringMap(v) }
   if v, ok := d.GetOk("schedule_info"); ok { req.ScheduleInfo = toStringMap(v) }
   if v, ok := d.GetOk("minutes"); ok { val := v.(int); req.Minutes = &val }
   if v, ok := d.GetOk("hours"); ok { val := v.(int); req.Hours = &val }
   if v, ok := d.GetOk("paused"); ok { val := v.(bool); req.Paused = &val }
   if v, ok := d.GetOk("recurring"); ok { val := v.(bool); req.Recurring = &val }
   if v, ok := d.GetOk("run_next_interval"); ok { val := v.(bool); req.RunNextInterval = &val }
   if v, ok := d.GetOk("schedule_name"); ok { val := v.(string); req.ScheduleName = &val }
   if v, ok := d.GetOk("schedule_type"); ok { val := v.(string); req.ScheduleType = &val }
   if v, ok := d.GetOk("starting_date"); ok { val := v.(string); req.StartingDate = &val }
   if v, ok := d.GetOk("monthdays"); ok { req.Monthdays = toIntSlice(v) }
   if v, ok := d.GetOk("weekdays"); ok { req.Weekdays = toWeekdaysSlice(v) }
   _, err = client.ModifySchedule(req)
   if err != nil {
      return err
   }
   return resourceElementswScheduleRead(d, m)
}

// resourceElementswScheduleDelete handles the Delete operation
func resourceElementswScheduleDelete(d *schema.ResourceData, m interface{}) error {
   client := m.(*Client)
   id, err := toIntID(d.Id())
   if err != nil {
      return err
   }
   toDel := true
   req := ModifyScheduleRequest{ScheduleID: id, ToBeDeleted: &toDel}
   _, err = client.ModifySchedule(req)
   if err != nil {
      return err
   }
   d.SetId("")
   return nil
}

// helper: convert interface{} to map[string]interface{}
func toStringMap(v interface{}) map[string]interface{} {
   if v == nil {
      return nil
   }
   if m, ok := v.(map[string]interface{}); ok {
      return m
   }
   return nil
}

// helper: convert interface{} list to []int

// helper: convert interface{} list to []map[string]int
func toWeekdaysSlice(v interface{}) []map[string]int {
   if v == nil {
      return nil
   }
   arr, ok := v.([]interface{})
   if !ok {
      return nil
   }
   out := make([]map[string]int, len(arr))
   for i, x := range arr {
      out[i], _ = x.(map[string]int)
   }
   return out
}

// helper: parse ID from string to int
func toIntID(id string) (int, error) {
   var i int
   _, err := fmt.Sscanf(id, "%d", &i)
   return i, err
}

