# elementsw_volume_iqn (Data Source)

Use this data source to construct the full iSCSI IQN and target portal for a specific SolidFire volume. This is useful for passing connection details to other providers (like `libvirt` or `vault`) without performing an API lookup if the basic identifiers are already known.

## Example Usage

```hcl
data "elementsw_cluster" "current" {}

resource "elementsw_volume" "disk" {
  name       = "my-disk"
  total_size = 1000000000
}

data "elementsw_volume_iqn" "disk_conn" {
  unique_id = data.elementsw_cluster.current.unique_id
  name      = elementsw_volume.disk.name
  volume_id = elementsw_volume.disk.id
  svip      = data.elementsw_cluster.current.svip
}

output "iscsi_target" {
  value = data.elementsw_volume_iqn.disk_conn.iqn
}
```

## Argument Reference

The following arguments are supported:

* `unique_id` - (Required) The Unique ID of the cluster (from `elementsw_cluster` data source).
* `name` - (Required) The name of the volume.
* `volume_id` - (Required) The numeric ID of the volume.
* `svip` - (Required) The Storage Virtual IP (SVIP) of the cluster.

## Attribute Reference

The following attributes are exported:

* `iqn` - The formatted iSCSI IQN.
* `target_portal` - The SVIP with the default iSCSI port (e.g., `192.168.1.34:3260`).
