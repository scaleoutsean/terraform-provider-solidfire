[![Go](https://github.com/scaleoutsean/terraform-provider-solidfire/actions/workflows/go.yml/badge.svg)](https://github.com/scaleoutsean/terraform-provider-solidfire/actions/workflows/go.yml)

# Terraform Provider for SolidFire

This is the Terraform Provider for SolidFire, used to configure resources on NetApp HCI or NetApp SolidFire storage clusters via the Element API.

For general information about Terraform, visit the [official website](https://terraform.io/).

**NOTE:** This repository is a community-maintained fork and is not associated with NetApp. To help with disambiguation, "SolidFire" is used in place of "ElementSW". This provider is tested with SolidFire version 12.5. Newer versions that use the 12.5 API endpoint should perfectly.

## Using the Provider

The provider is available on the [Terraform Registry](https://registry.terraform.io/providers/scaleoutsean/solidfire/latest).

To use the provider, require it in your Terraform configuration:

```hcl
terraform {
  required_providers {
    solidfire = {
      source  = "scaleoutsean/solidfire"
      version = "~> 0.4.2" # replace with the latest version
    }
  }
}

provider "solidfire" {
  # configuration options...
}
```

## Naming Conventions

SolidFire does not require all resource names to be unique; they are internally treated as labels while resources are uniquely identified by IDs. However, these IDs are generated on the fly and are not user-friendly.

This provider assumes that resource names are unique and enforces this within its scope. This functions correctly if everything is managed through Terraform, but could raise conflicts if the rule is violated by managing resources directly via the SolidFire UI or API outside of Terraform.

**WARNING:** Be careful with **immutable** resource properties (e.g., volume name and IQN name). Terraform will not be able to update these in-place. If changed, Terraform will destroy the existing resource and create a new one.

## Developing & Testing

If you wish to contribute to the provider, you'll need [Go](https://golang.org/) installed (see `go.mod` for the required version).

1. Clone the repository.
2. Run `make build` to compile the provider plugin.
3. Set the required environment variables, including `TF_ACC=1` to enable Terraform Acceptance Tests, and `SOLIDFIRE_ACC=1` for provider-specific protections. 
4. Run `go test ./... -v -timeout 15m` to run acceptance tests against a live SolidFire cluster.

```sh
export SOLIDFIRE_SERVER="elementsw.cluster.ipv4"
export SOLIDFIRE_USERNAME="admin"
export SOLIDFIRE_PASSWORD="changeme"
export SOLIDFIRE_API_VERSION="12.5"
export TF_ACC="1"
# Some provider tests require this additional flag to run safely
export SOLIDFIRE_ACC="1" 
go test ./... -v -timeout 15m
```

To test two clusters (cluster and volume pairing), provide environment variables for the second cluster:

```sh
export SOLIDFIRE_SERVER_DR=...
export SOLIDFIRE_USERNAME_DR=...
export SOLIDFIRE_PASSWORD_DR=...
export SOLIDFIRE_API_VERSION_DR="12.5"
```

Please submit issues or pull requests to the [GitHub repository](https://github.com/scaleoutsean/terraform-provider-solidfire).
