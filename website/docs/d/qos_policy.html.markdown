---
layout: "elementsw"
page_title: "ElementSW: elementsw_qos_policy"
sidebar_current: "docs-elementsw-datasource-qos-policy"
description: |-
  Get information about a QoS policy on a SolidFire cluster.
---

# elementsw_qos_policy

Use this data source to get information about a QoS policy on a SolidFire cluster.

## Example Usage

```hcl
data "elementsw_qos_policy" "standard" {
  name = "Standard"
}

resource "elementsw_volume" "example" {
  name       = "example-vol"
  account_id = elementsw_account.example.account_id
  total_size = 1000000000
  enable512e = true
  qos_policy_id = data.elementsw_qos_policy.standard.qos_policy_id
}
```

## Argument Reference

The following arguments are supported:

* `qos_policy_id` - (Optional) The ID of the QoS policy to look up.
* `name` - (Optional) The name of the QoS policy to look up.

One of `qos_policy_id` or `name` must be specified.

## Attributes Reference

The following attributes are exported:

* `qos_policy_id` - The ID of the QoS policy.
* `name` - The name of the QoS policy.
* `min_iops` - The minimum IOPS of the QoS policy.
* `max_iops` - The maximum IOPS of the QoS policy.
* `burst_iops` - The burst IOPS of the QoS policy.
