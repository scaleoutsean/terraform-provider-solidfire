resource "solidfire_volume" "volume" {
  name       = "my-volume"
  account    = solidfire_account.k8s_account.id
  total_size = 1073742000
  enable512e = false
  min_iops   = 100
  max_iops   = 150
  burst_iops = 200
}
