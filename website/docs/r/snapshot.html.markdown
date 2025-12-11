---
title: "elementsw_snapshot"
---

# Resource: elementsw_snapshot

Manages a SolidFire snapshot (individual or group).

## Example Usage

```hcl
resource "elementsw_snapshot" "example" {
  volume_id = 1
  name = "my-snapshot"
}

resource "elementsw_snapshot" "group" {
  volume_ids = [1,2]
  name = "my-group-snapshot"
}
```

## Argument Reference

- `volume_id` (Optional) - The volume ID for an individual snapshot.
- `volume_ids` (Optional) - List of volume IDs for a group snapshot.
- `name` (Optional) - The snapshot name.
- `snapmirror_label` (Optional) - SnapMirror label (ONTAP-related).
- `enable_remote_replication` (Optional) - Enable remote replication.
- `ensure_serial_creation` (Optional) - Ensure serial creation for group snapshots.
- `retention` (Optional) - Retention period.
- `expiration_time` (Optional) - Expiration time.
- `attributes` (Optional) - Map of custom attributes.
- `save_members` (Optional, Delete only) - Save member snapshots when deleting a group snapshot.

## Attribute Reference

- `created_snapshot_id` - The created snapshot ID.
- `created_group_snapshot_id` - The created group snapshot ID.
- `created_group_snapshot_uuid` - The created group snapshot UUID.
- `create_time` - The snapshot creation time.
