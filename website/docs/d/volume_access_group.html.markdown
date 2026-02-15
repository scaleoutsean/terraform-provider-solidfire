---
layout: "elementsw"
page_title: "ElementSW: elementsw_volume_access_group"
sidebar_current: "docs-elementsw-datasource-volume-access-group"
description: |-
  Get information about a Volume Access Group on a SolidFire cluster.
---

# elementsw_volume_access_group

Use this data source to get information about a Volume Access Group (VAG) on a SolidFire cluster.

## Example Usage

```hcl
data "elementsw_volume_access_group" "standard" {
  name = "standard-vag"
}

output "vag_id" {
  value = data.elementsw_volume_access_group.standard.volume_access_group_id
}
```

## Argument Reference

The following arguments are supported:

* `volume_access_group_id` - (Optional) The ID of the VAG to look up.
* `name` - (Optional) The name of the VAG to look up.

One of `volume_access_group_id` or `name` must be specified.

## Attributes Reference

The following attributes are exported:

* `volume_access_group_id` - The ID of the VAG.
* `name` - The name of the VAG.
* `initiators` - The list of initiators (IQNs) in the VAG.
* `volumes` - The list of volume IDs in the VAG.
