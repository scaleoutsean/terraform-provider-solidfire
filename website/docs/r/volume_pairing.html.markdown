---
title: "elementsw_volume_pairing"
---

# Resource: elementsw_volume_pairing

Manages volume-to-volume replication pairing.

## Example Usage

```hcl
resource "elementsw_volume" "primary" {
  name       = "primary-vol"
  total_size = 1000000000
}

resource "elementsw_volume_pairing" "dr_pair" {
  volume_id = elementsw_volume.primary.id
  target_cluster {
    endpoint = "https://10.20.20.20/json-rpc/12.5"
    username = "admin"
    password = "password"
  }
}
```

## Argument Reference

- `volume_id` (Required) - The ID of the primary volume to be paired.
- `target_cluster` (Optional) - Connection info for the destination cluster. If provided, the provider will attempt to find a volume with the same name on the target cluster and complete the pairing automatically.
  - `endpoint` (Required) - API endpoint.
  - `username` (Required) - Admin username.
  - `password" (Required) - Admin password.
- `pairing_key` (Optional/Computed) - The pairing key generated for this volume pair.

## Attribute Reference

- `pairing_key` - The key used to pair with the remote volume.
