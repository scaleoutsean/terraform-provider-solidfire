---
title: "elementsw_schedule"
---

# Resource: elementsw_schedule

Manages a SolidFire snapshot schedule.

## Example Usage

```hcl
resource "elementsw_schedule" "example" {
  name = "daily-snap"
  volume_id = 1
  hours = [0]
  minutes = [0]
  attributes = {
    key = "value"
  }
}
```

## Argument Reference

- `name` (Required) - The schedule name.
- `volume_id` (Required) - The volume ID to snapshot.
- `hours` (Optional) - List of hours to run.
- `minutes` (Optional) - List of minutes to run.
- `attributes` (Optional) - Map of custom attributes.

## Attribute Reference

- `schedule_id` - The SolidFire schedule ID.
