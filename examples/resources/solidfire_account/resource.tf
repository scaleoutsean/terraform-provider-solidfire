resource "solidfire_account" "k8s_account" {
  username         = "k8s-cluster"
  target_secret    = "targetsecret123"
  initiator_secret = "initsecret123"
}
