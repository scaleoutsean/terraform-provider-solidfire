# elementsw_volume (Data Source)

Use this data source to get information about an existing volume on a SolidFire cluster.

## Example Usage

```hcl
data "elementsw_volume" "example" {
  name = "my-volume"
}

output "volume_id" {
  value = data.elementsw_volume.example.volume_id
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` - (Optional) The ID of the volume to look up.
* `name` - (Optional) The name of the volume to look up.

Either `volume_id` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `account_id` - The ID of the account that owns the volume.
* `total_size` - The total size of the volume in bytes.
* `iqn` - The iSCSI Qualified Name (IQN) of the volume.
* `access` - The access mode of the volume (e.g., `readWrite`, `replicationTarget`).
