resource "solidfire_qos_policy" "high_priority" {
  name       = "HighPriority"
  min_iops   = 1000
  max_iops   = 5000
  burst_iops = 10000
}
