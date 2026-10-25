# CHANGELOG

## v0.5.0 (2026/10/25)

* Updates: Update direct dependency SolidFire Go v1.0.7 (updated to Go 1.27.3)

## v0.4.9 (2026/10/25)

* Updates: Go 1.27.3 and various other 3rd party dependencies

## v0.4.8 (2026/10/25)

* Update gRPC and indirect deps

## v0.4.7 (2026/07/14)

* Update Go to 1.26.5 (CVE), SolidFire Go SDK to 1.0.6 (due to Go CVE) and other deps

## v0.4.6 (2026/05/16)

* Update golang.org/x/net@v0.53.0 (GO-2026-4918) and SolidFire Go SDK

## v0.4.5 (2026/05/16)

* Update documentation

## v0.4.4 (2026/04/05)

* Restructure provider directory layout (`elementsw` to `solidfire`) and generate new docs
* **Fix**: Remove `burst_time` from `solidfire_qos_policy` resource and tests (unsupported/invalid parameter)
* **Fix**: Ensure mutually exclusive configuration of `account` and `account_id` in `solidfire_volume` via `ExactlyOneOf`
* Dependencies update

## v0.4.3 (2026/03/06)

* Minor documentation updates

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
