# elementsw_cluster (Data Source)

Use this data source to get information about the connected SolidFire cluster.

## Example Usage

```hcl
data "elementsw_cluster" "current" {}

output "cluster_version" {
  value = data.elementsw_cluster.current.cluster_version
}
```

## Attribute Reference

The following attributes are exported:

* `name` - The name of the cluster.
* `unique_id` - The unique ID assigned to the cluster.
* `cluster_version` - The software version running on the cluster.
* `cluster_api_version` - The API version supported by the cluster.
* `mvip` - The Management Virtual IP address.
* `svip` - The Storage Virtual IP address.
