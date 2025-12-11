---
title: "elementsw_volumes_by_account"
---

# Data Source: elementsw_volumes_by_account

Use this data source to get a list of SolidFire volume IDs for a given account.

## Example Usage

```hcl
data "elementsw_volumes_by_account" "example" {
  account_id = 123
}

output "volume_ids" {
  value = data.elementsw_volumes_by_account.example.volume_ids
}
```

## Argument Reference

- `account_id` (Required) - The SolidFire account ID to query.

## Attribute Reference

- `volume_ids` - List of volume IDs belonging to the account.
