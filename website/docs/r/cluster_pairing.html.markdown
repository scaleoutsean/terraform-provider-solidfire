---
title: "elementsw_cluster_pairing"
---

# Resource: elementsw_cluster_pairing

Manages a cluster-to-cluster trust relationship (pairing) for replication.

## Example Usage (Automated Pairing)

This workflow uses one provider instance (the "local" one) and connects to a remote cluster to perform the bidirectional exchange automatically.

```hcl
resource "elementsw_cluster_pairing" "dr_pair" {
  source_cluster {
    endpoint = "https://10.10.10.10/json-rpc/12.5"
    username = "admin"
    password = "password"
  }
  target_cluster {
    endpoint = "https://10.20.20.20/json-rpc/12.5"
    username = "admin"
    password = "password"
  }
}
```

## Example Usage (Manual Key Exchange)

If you already have a pairing key from a `StartClusterPairing` operation, you can provide it directly to the target cluster.

```hcl
resource "elementsw_cluster_pairing" "dr_pair" {
  pairing_key = "ey..." # Value from StartClusterPairing
  target_cluster {
    endpoint = "https://10.20.20.20/json-rpc/12.5"
    username = "admin"
    password = "password"
  }
}
```

## Argument Reference

- `target_cluster` (Required) - Connection info for the cluster where pairing will be completed.
  - `endpoint` (Required) - API endpoint (e.g., https://10.1.1.1/json-rpc/12.5).
  - `username` (Required) - Admin username.
  - `password" (Required) - Admin password.
- `source_cluster` (Optional) - Connection info for the cluster where pairing will be started. Required for automated workflow.
- `pairing_key` (Optional) - Manual pairing key.

## Attribute Reference

- `cluster_pair_id` - The ID of the cluster pair.
- `cluster_name` - Name of the remote cluster.
- `status` - Status of the pairing.
