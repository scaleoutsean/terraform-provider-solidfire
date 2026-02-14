package elementsw

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

// resourceElementswSnapshot returns the Terraform resource for SolidFire snapshots (individual or group)
func resourceElementswSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementswSnapshotCreate,
		Read:   resourceElementswSnapshotRead,
		Update: resourceElementswSnapshotUpdate,
		Delete: resourceElementswSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"volume_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"snapshot_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"group_snapshot_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapmirror_label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_remote_replication": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ensure_serial_creation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"retention": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"save_members": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// Output fields
			"created_snapshot_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_group_snapshot_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceElementswSnapshotCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	// If group_snapshot, use CreateGroupSnapshot
	if _, ok := d.GetOk("volume_ids"); ok {
		// Create group snapshot
		volumes := []int64{}
		if v, ok := d.GetOk("volume_ids"); ok {
			for _, id := range v.([]interface{}) {
				volumes = append(volumes, int64(id.(int)))
			}
		}
		req := sdk.CreateGroupSnapshotRequest{
			Volumes: volumes,
			Name:    d.Get("name").(string),
		}
		if v, ok := d.GetOk("enable_remote_replication"); ok {
			req.EnableRemoteReplication = v.(bool)
		}
		if v, ok := d.GetOk("retention"); ok {
			req.Retention = v.(string)
		}
		// ExpirationTime and EnsureSerialCreation are not in CreateGroupSnapshotRequest in this SDK version
		resp, err := client.CreateGroupSnapshot(&req)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("group-%d", resp.GroupSnapshotID))
		d.Set("created_group_snapshot_id", int(resp.GroupSnapshotID))
		return resourceElementswSnapshotRead(d, m)
	}
	// Otherwise, create individual snapshot
	req := sdk.CreateSnapshotRequest{
		VolumeID:        int64(d.Get("volume_id").(int)),
		Name:            d.Get("name").(string),
		SnapMirrorLabel: d.Get("snapmirror_label").(string),
	}
	if v, ok := d.GetOk("enable_remote_replication"); ok {
		req.EnableRemoteReplication = v.(bool)
	}
	if v, ok := d.GetOk("retention"); ok {
		req.Retention = v.(string)
	}
	// ExpirationTime and EnsureSerialCreation are not in CreateSnapshotRequest in this SDK version
	resp, err := client.CreateSnapshot(&req)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("snap-%d", resp.SnapshotID))
	d.Set("created_snapshot_id", int(resp.SnapshotID))
	d.Set("create_time", resp.Snapshot.CreateTime)
	return resourceElementswSnapshotRead(d, m)
}

func resourceElementswSnapshotUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	idStr := d.Id()
	if strings.HasPrefix(idStr, "group-") {
		id, _ := strconv.ParseInt(strings.TrimPrefix(idStr, "group-"), 10, 64)
		req := sdk.ModifyGroupSnapshotRequest{
			GroupSnapshotID: id,
		}
		if v, ok := d.GetOk("enable_remote_replication"); ok {
			req.EnableRemoteReplication = v.(bool)
		}
		if v, ok := d.GetOk("expiration_time"); ok {
			req.ExpirationTime = v.(string)
		}
		if v, ok := d.GetOk("snapmirror_label"); ok {
			req.SnapMirrorLabel = v.(string)
		}
		// Note: Name is not in ModifyGroupSnapshotRequest in this SDK version
		return client.ModifyGroupSnapshot(&req)
	} else if strings.HasPrefix(idStr, "snap-") {
		id, _ := strconv.ParseInt(strings.TrimPrefix(idStr, "snap-"), 10, 64)
		req := sdk.ModifySnapshotRequest{
			SnapshotID: id,
		}
		if v, ok := d.GetOk("enable_remote_replication"); ok {
			req.EnableRemoteReplication = v.(bool)
		}
		if v, ok := d.GetOk("expiration_time"); ok {
			req.ExpirationTime = v.(string)
		}
		if v, ok := d.GetOk("snapmirror_label"); ok {
			req.SnapMirrorLabel = v.(string)
		}
		// Note: Name is not in ModifySnapshotRequest in this SDK version
		return client.ModifySnapshot(&req)
	}
	return nil
}

func resourceElementswSnapshotDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	idStr := d.Id()
	if strings.HasPrefix(idStr, "group-") {
		id, _ := strconv.ParseInt(strings.TrimPrefix(idStr, "group-"), 10, 64)
		saveMembers := false
		if v, ok := d.GetOk("save_members"); ok {
			saveMembers = v.(bool)
		}
		return client.DeleteGroupSnapshot(id, saveMembers)
	} else if strings.HasPrefix(idStr, "snap-") {
		id, _ := strconv.ParseInt(strings.TrimPrefix(idStr, "snap-"), 10, 64)
		return client.DeleteSnapshot(id)
	}
	return nil
}

func resourceElementswSnapshotRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	idStr := d.Id()

	if strings.HasPrefix(idStr, "group-") {
		id, _ := strconv.ParseInt(strings.TrimPrefix(idStr, "group-"), 10, 64)
		vols := []int64{}
		if v, ok := d.GetOk("volume_ids"); ok {
			for _, vid := range v.([]interface{}) {
				vols = append(vols, int64(vid.(int)))
			}
		}
		res, err := client.ListGroupSnapshots(vols)
		if err != nil {
			return err
		}
		for _, gs := range res {
			if gs.GroupSnapshotID == id {
				d.Set("name", gs.Name)
				d.Set("create_time", gs.CreateTime)
				return nil
			}
		}
		d.SetId("") // Not found
	} else if strings.HasPrefix(idStr, "snap-") {
		id, _ := strconv.ParseInt(strings.TrimPrefix(idStr, "snap-"), 10, 64)
		volID := int64(d.Get("volume_id").(int))
		res, err := client.ListSnapshots(volID)
		if err != nil {
			return err
		}
		for _, s := range res {
			if s.SnapshotID == id {
				d.Set("name", s.Name)
				d.Set("create_time", s.CreateTime)
				d.Set("snapmirror_label", s.SnapMirrorLabel)
				return nil
			}
		}
		d.SetId("") // Not found
	}

	return nil
}
