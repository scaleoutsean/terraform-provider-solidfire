# CHANGELOG 

## v0.3.0 (2026/02/23)

BREAKING CHANGES:

* **Renamed Go module:** `github.com/netapp/terraform-provider-netapp-elementsw` -> `github.com/scaleoutsean/terraform-provider-solidfire`
* **SDK Replacement:** Replaced `solidfire-sdk-go` with `github.com/scaleoutsean/solidfire-go`

FEATURES:

* **New Resource:** `elementsw_cluster_pairing`
* **New Resource:** `elementsw_schedule`
* **New Resource:** `elementsw_snapshot`
* **New Resource:** `elementsw_volume_pairing`
* **New Data Source:** `elementsw_account`
* **New Data Source:** `elementsw_cluster`
* **New Data Source:** `elementsw_cluster_stats`
* **New Data Source:** `elementsw_initiator`
* **New Data Source:** `elementsw_qos_policy`
* **New Data Source:** `elementsw_volume`
* **New Data Source:** `elementsw_volume_access_group`
* **New Data Source:** `elementsw_volume_iqn`
* **New Data Source:** `elementsw_volumes_by_account`

IMPROVEMENTS:

* Integrated `terraform-plugin-testing` framework
* Updated Go version and dependencies

## 0.2.1 (2025-07-25)

* **New Resource:** QoS Policy (List, Get) in `elementsw_qos_policy`

## 0.2.0 (Unreleased)

* **New:** Forked from Terraform NetApp ElementSW Provider (repository)
* **New:** Updated to work with Terraform 1.5 and Terraform Plugin SDK v2.30

## 0.1.0 (Unreleased)

FEATURES:

* **New Resource:** `elementsw_initiator`
* **New Resource:** `elementsw_volume`
* **New Resource:** `elementsw_volume_access_group`
* **New Resource:** `elementsw_account`
