# elementsw_cluster_stats (Data Source)

Use this data source to get real-time performance and capacity metrics from a SolidFire cluster.

## Example Usage

```hcl
data "elementsw_cluster_stats" "current" {}

output "iops" {
  value = data.elementsw_cluster_stats.current.metrics[0].actual_iops
}

output "utilization" {
  value = data.elementsw_cluster_stats.current.metrics[0].cluster_utilization
}
```

## Attribute Reference

The following attributes are exported:

### Summary Stats
* `volume_count` - Total number of volumes in the cluster.
* `node_count` - Total number of active nodes in the cluster.
* `volumes_per_node` - Average number of volumes per node.
* `compression_factor` - Efficiency ratio (calculated as Provisioned Space / Used Space).

### `metrics` Block
* `actual_iops` - The current actual IOPS for the entire cluster.
* `average_iop_size` - Average size in bytes of I/O operations.
* `cluster_utilization` - Percentage of the cluster's performance capacity currently being used.
* `latency_usec` - Average cluster latency in microseconds.
* `read_bytes` - Total bytes read per second.
* `read_ops` - Total read operations per second.
* `write_bytes` - Total bytes written per second.
* `write_ops` - Total write operations per second.
* `timestamp` - The ISO 8601 timestamp of the data.

### `capacity` Block
* `active_block_space` - Space in bytes used by active blocks.
* `max_iops` - Maximum IOPS capacity of the cluster.
* `max_used_space` - Total capacity of the cluster in bytes.
* `provisioned_space` - Total provisioned space across all volumes.
* `used_space` - Total used space on the cluster.
* `unique_blocks` - Space used by unique blocks (deduplicated).
* `zero_blocks` - Space that would have been used by zero blocks.
