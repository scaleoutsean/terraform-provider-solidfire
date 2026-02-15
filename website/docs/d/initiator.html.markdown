---
layout: "elementsw"
page_title: "ElementSW: elementsw_initiator"
sidebar_current: "docs-elementsw-datasource-initiator"
description: |-
  Get information about an iSCSI initiator on a SolidFire cluster.
---

# elementsw_initiator

Use this data source to get information about an iSCSI initiator on a SolidFire cluster.

## Example Usage

```hcl
data "elementsw_initiator" "node1" {
  name = "iqn.1994-05.com.redhat:node1"
}

output "initiator_id" {
  value = data.elementsw_initiator.node1.initiator_id
}
```

## Argument Reference

The following arguments are supported:

* `initiator_id` - (Optional) The ID of the initiator to look up.
* `name` - (Optional) The name (IQN) of the initiator to look up.

One of `initiator_id` or `name` must be specified.

## Attributes Reference

The following attributes are exported:

* `initiator_id` - The ID of the initiator.
* `name` - The name (IQN) of the initiator.
* `alias` - The alias for the initiator.
* `volume_access_group_ids` - The list of Volume Access Group IDs this initiator belongs to.
