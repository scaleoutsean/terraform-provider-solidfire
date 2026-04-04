resource "solidfire_initiator" "test_initiator" {
  name                   = "iqn.1993-08.org.debian:01:my-initiator"
  alias                  = "my-initiator"
  volume_access_group_id = solidfire_volume_access_group.test_group.id
}
