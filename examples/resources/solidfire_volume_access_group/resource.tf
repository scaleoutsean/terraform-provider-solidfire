resource "solidfire_volume_access_group" "test-group" {
  name     = "my-vag"
  volumes  = [solidfire_volume.volume.id]
}
