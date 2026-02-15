---
layout: "elementsw"
page_title: "ElementSW: elementsw_account"
sidebar_current: "docs-elementsw-datasource-account"
description: |-
  Get information about a tenant account on a SolidFire cluster.
---

# elementsw_account

Use this data source to get information about a tenant account on a SolidFire cluster.

## Example Usage

```hcl
data "elementsw_account" "example" {
  username = "example-user"
}

output "account_id" {
  value = data.elementsw_account.example.account_id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The ID of the account to look up.
* `username` - (Optional) The username of the account to look up.

One of `account_id` or `username` must be specified.

## Attributes Reference

The following attributes are exported:

* `account_id` - The ID of the account.
* `username` - The username of the account.
* `status` - The status of the account.
* `initiator_secret` - (Sensitive) The initiator CHAP secret. (Currently dropped for security reasons).
* `target_secret` - (Sensitive) The target CHAP secret. (Currently dropped for security reasons).
