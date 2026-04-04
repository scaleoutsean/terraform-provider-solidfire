resource "solidfire_schedule" "daily" {
  name      = "daily-snapshot"
  volume_id = solidfire_volume.volume.id
  minutes   = 0
  hours     = 2
}
