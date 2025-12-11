---
title: "elementsw_qos_policy"
---

# Resource: elementsw_qos_policy

Manages a SolidFire QoS policy.

## Example Usage

```hcl
resource "elementsw_qos_policy" "example" {
  name = "mypolicy"
  min_iops = 1000
  max_iops = 4000
  burst_iops = 6000
}
```

## Argument Reference

- `name` (Required) - The policy name.
- `min_iops` (Required) - Minimum IOPS.
- `max_iops` (Required) - Maximum IOPS.
- `burst_iops` (Required) - Burst IOPS.

## Attribute Reference

- `qos_policy_id` - The SolidFire QoS policy ID.
