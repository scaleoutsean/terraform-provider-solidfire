resource "solidfire_snapshot" "snap1" {
  name      = "my-snapshot"
  volume_id = solidfire_volume.volume.id
}
