package elementsw

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceElementswSnapshot returns the Terraform resource for SolidFire snapshots (individual or group)
func resourceElementswSnapshot() *schema.Resource {
   return &schema.Resource{
	  Create: resourceElementswSnapshotCreate,
	  Read:   resourceElementswSnapshotRead,
	  Update: resourceElementswSnapshotUpdate,
	  Delete: resourceElementswSnapshotDelete,
	  Schema: map[string]*schema.Schema{
		 "volume_id": {Type: schema.TypeInt, Optional: true},
		 "snapshot_id": {Type: schema.TypeInt, Optional: true},
		 "group_snapshot_id": {Type: schema.TypeInt, Optional: true},
		 "name": {Type: schema.TypeString, Optional: true},
		 "snapmirror_label": {Type: schema.TypeString, Optional: true},
		 "enable_remote_replication": {Type: schema.TypeBool, Optional: true},
		 "ensure_serial_creation": {Type: schema.TypeBool, Optional: true},
		 "retention": {Type: schema.TypeString, Optional: true},
		 "expiration_time": {Type: schema.TypeString, Optional: true},
		 "attributes": {Type: schema.TypeMap, Optional: true},
		 "save_members": {Type: schema.TypeBool, Optional: true},
		 // Output fields
		 "created_snapshot_id": {Type: schema.TypeInt, Computed: true},
		 "created_group_snapshot_id": {Type: schema.TypeInt, Computed: true},
		 "created_group_snapshot_uuid": {Type: schema.TypeString, Computed: true},
		 "create_time": {Type: schema.TypeString, Computed: true},
	  },
   }
}

func resourceElementswSnapshotCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	// If group_snapshot, use CreateGroupSnapshot
	if v, ok := d.GetOk("group_snapshot_id"); ok && v.(int) == 0 {
		// Create group snapshot
		volumes := []int{}
		if v, ok := d.GetOk("volume_ids"); ok {
			volumes = toIntSlice(v)
		}
		req := CreateGroupSnapshotRequest{
			Volumes: volumes,
			Name: d.Get("name").(string),
		}
		if v, ok := d.GetOk("enable_remote_replication"); ok { b := v.(bool); req.EnableRemoteReplication = &b }
		if v, ok := d.GetOk("ensure_serial_creation"); ok { b := v.(bool); req.EnsureSerialCreation = &b }
		if v, ok := d.GetOk("retention"); ok { req.Retention = v.(string) }
		if v, ok := d.GetOk("expiration_time"); ok { req.ExpirationTime = v.(string) }
		if v, ok := d.GetOk("attributes"); ok { req.Attributes = v.(map[string]interface{}) }
		resp, err := client.CreateGroupSnapshot(req)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("group-%d", resp.GroupSnapshotID))
		d.Set("created_group_snapshot_id", resp.GroupSnapshotID)
		d.Set("created_group_snapshot_uuid", resp.GroupSnapshotUUID)
		return resourceElementswSnapshotRead(d, m)
	}
	// Otherwise, create individual snapshot
	req := CreateSnapshotRequest{
		VolumeID: d.Get("volume_id").(int),
		Name: d.Get("name").(string),
		SnapMirrorLabel: d.Get("snapmirror_label").(string),
	}
	if v, ok := d.GetOk("snapshot_id"); ok { id := v.(int); req.SnapshotID = &id }
	if v, ok := d.GetOk("enable_remote_replication"); ok { b := v.(bool); req.EnableRemoteReplication = &b }
	if v, ok := d.GetOk("ensure_serial_creation"); ok { b := v.(bool); req.EnsureSerialCreation = &b }
	if v, ok := d.GetOk("retention"); ok { req.Retention = v.(string) }
	if v, ok := d.GetOk("expiration_time"); ok { req.ExpirationTime = v.(string) }
	if v, ok := d.GetOk("attributes"); ok { req.Attributes = v.(map[string]interface{}) }
	resp, err := client.CreateSnapshot(req)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("snap-%d", resp.Snapshot.SnapshotID))
	d.Set("created_snapshot_id", resp.Snapshot.SnapshotID)
	d.Set("create_time", resp.Snapshot.CreateTime)
	return resourceElementswSnapshotRead(d, m)
}

func resourceElementswSnapshotUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if v, ok := d.GetOk("group_snapshot_id"); ok && v.(int) != 0 {
		req := ModifyGroupSnapshotRequest{
			GroupSnapshotID: v.(int),
			Name: d.Get("name").(string),
		}
		if v, ok := d.GetOk("enable_remote_replication"); ok { b := v.(bool); req.EnableRemoteReplication = &b }
		if v, ok := d.GetOk("expiration_time"); ok { req.ExpirationTime = v.(string) }
		return client.ModifyGroupSnapshot(req)
	}
	if v, ok := d.GetOk("snapshot_id"); ok && v.(int) != 0 {
		req := ModifySnapshotRequest{
			SnapshotID: v.(int),
			Name: d.Get("name").(string),
			SnapMirrorLabel: d.Get("snapmirror_label").(string),
		}
		if v, ok := d.GetOk("enable_remote_replication"); ok { b := v.(bool); req.EnableRemoteReplication = &b }
		if v, ok := d.GetOk("expiration_time"); ok { req.ExpirationTime = v.(string) }
		return client.ModifySnapshot(req)
	}
	return nil
}

func resourceElementswSnapshotDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if v, ok := d.GetOk("group_snapshot_id"); ok && v.(int) != 0 {
		saveMembers := false
		if v, ok := d.GetOk("save_members"); ok { saveMembers = v.(bool) }
		req := DeleteGroupSnapshotRequest{
			GroupSnapshotID: v.(int),
			SaveMembers: saveMembers,
		}
		return client.DeleteGroupSnapshot(req)
	}
	if v, ok := d.GetOk("snapshot_id"); ok && v.(int) != 0 {
		req := DeleteSnapshotRequest{SnapshotID: v.(int)}
		return client.DeleteSnapshot(req)
	}
	return nil
}

// resourceElementswSnapshotRead handles the Read operation for snapshots
func resourceElementswSnapshotRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	var groupSnapshots []map[string]interface{}
	var snapshots []map[string]interface{}

	if v, ok := d.GetOk("volume_ids"); ok {
		volIDs := toIntSlice(v)
		if len(volIDs) > 0 {
			// ListGroupSnapshots for specified volumes
			res, err := client.ListGroupSnapshots(ListGroupSnapshotsRequest{Volumes: volIDs})
			if err != nil {
				return err
			}
			for _, gs := range res.GroupSnapshots {
				members := make([]map[string]interface{}, len(gs.Members))
				for i, m := range gs.Members {
					members[i] = map[string]interface{}{
						"snapshot_id": m.SnapshotID,
						"volume_id": m.VolumeID,
						"create_time": m.CreateTime,
					}
				}
				groupSnapshots = append(groupSnapshots, map[string]interface{}{
					"group_snapshot_id": gs.GroupSnapshotID,
					"group_snapshot_uuid": gs.GroupSnapshotUUID,
					"members": members,
				})
			}
		}
	}
	// Always list individual snapshots
	res, err := client.ListSnapshots(ListSnapshotsRequest{})
	if err != nil {
		return err
	}
	for _, s := range res.Snapshots {
		snapshots = append(snapshots, map[string]interface{}{
			"snapshot_id": s.SnapshotID,
			"volume_id": s.VolumeID,
			"create_time": s.CreateTime,
		})
	}
	// Set results
	if len(groupSnapshots) > 0 {
		d.Set("group_snapshots", groupSnapshots)
	}
	d.Set("snapshots", snapshots)
	// Use a synthetic ID for the data source
	d.SetId(fmt.Sprintf("snapshots-%d", len(snapshots)+len(groupSnapshots)))
	return nil
}

// helper: convert interface{} list to []int
