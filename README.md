<!-- TOC -->

- [Terraform Provider SolidFire](#terraform-provider-solidfire)
  - [Naming Conventions](#naming-conventions)
  - [Using the Provider](#using-the-provider)
    - [Provider Documentation](#provider-documentation)
    - [Controlling the provider version](#controlling-the-provider-version)
  - [Building The Provider](#building-the-provider)
    - [Prerequisites](#prerequisites)
    - [Cloning the Project](#cloning-the-project)
    - [Running the Build](#running-the-build)
    - [Installing the Local Plugin](#installing-the-local-plugin)
  - [Developing the Provider](#developing-the-provider)
  - [Testing the Provider](#testing-the-provider)
    - [Configuring Environment Variables](#configuring-environment-variables)
      - [Using the `.tf-elementsw-devrc.mk` file](#using-the-tf-elementsw-devrcmk-file)
    - [Running the Acceptance Tests](#running-the-acceptance-tests)
    - [Walk-through example](#walk-through-example)
      - [Installing Go and Terraform](#installing-go-and-terraform)
      - [Installing dependencies](#installing-dependencies)
      - [Cloning the NetApp provider repository and building the provider](#cloning-the-netapp-provider-repository-and-building-the-provider)
      - [Sanity check](#sanity-check)

<!-- /TOC -->

# Terraform Provider SolidFire

This is the repository for the Terraform Provider (for) SolidFire, which can be used with Terraform to configure resources on NetApp HCI or NetApp SolidFire storage clusters.

For general information about Terraform, visit the [official website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform


**NOTE:** 

- This is a fork of NetApp-hosted "Terraform NetApp ElementSW Provider". 
- This repository is not associated with NetApp. To help with disambiguation, "NetApp" has been removed from the name, and "SolidFire" is used in place of "ElementSW". Terraform NetApp ElementSW Provider itself is based on a code initially developed by the SolidFire team for use with internal projects. The provider plugin was refactored to be published and maintained. It is possible that changes from Terraform SolidFire Provider may be submitted upstream to Terraform NetApp ElementSW Provider, but it's been forked because my primary goal is to experiment with it and not aim for pull request submission to Terraform NetApp ElementSW Provider

This provider was tested with SolidFire version 12.

## Naming Conventions

SolidFire does not require resource names to be unique. They are considered as 'labels' and resources in SolidFire are uniquely identified by IDs (integers). However, these ids are not user friendly, and as they are generated on the fly, they make it difficult to track resources and automate.

This provider assumes that resource names are unique, and enforces it within its scope. This is not an issue if everything is managed through Terraform, but could raise conflicts if the rule is not respected outside of Terraform.

## Using the Provider

The current version of this provider requires Terraform 1.5 to run.

**TODO** Download the provider from [Terraform Registry](https://registry.terraform.io/) if you don't want to build it from source. Note that you need to run `terraform init` to fetch the provider before deploying.

A how-to based on Terraform 1.5 and SolidFire 12.5 can be found in this repository.

### Provider Documentation

**TODO** The provider is documented [here](https://registry.terraform.io/providers/NetApp/netapp-elementsw/latest/docs).

Check the provider documentation for details on entering your connection information and how to get started with writing configuration SolidFire resources.

### Controlling the provider version

Note that you can also control the provider version. Since Terraform 0.13 this requires the use of a `required_providers` block in your Terraform configuration.

The syntax that loads the provider from Terraform Registry is as follows:

```hcl
required_providers {
  netapp-elementsw = {
    version = "~> 20.11"
    source  = "scaleoutsesan/solidfire"
  }
}
```

Version locking uses a pessimistic operator, so this version lock would mean anything within the 20.11 namespace, including or after 20.11.0. Read more [here][provider-vc] on provider version control.

For offline loading please see the walk-through for building from source.

[provider-vc]: https://www.terraform.io/docs/language/providers/requirements.html#version

## Building The Provider

This section is intended for developers.

### Prerequisites

If you wish to work on the provider, you'll first need [Go][go-website] installed on your machine (version 1.20+ (see `go.mod` for latest and accurate information) is **required** to build with current dependencies). You'll also need to correctly setup Go.

### Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/scaleoutsean/terraform-provider-solidfire`:

```sh
mkdir -p $GOPATH/src/github.com/netapp
cd $GOPATH/src/github.com/netapp
git clone https://github.com/NetApp/terraform-provider-solidfire.git
```

### Running the Build

After the clone has been completed, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/src/github.com/scaleoutsean/terraform-provider-solidfire
make build
```

### Installing the Local Plugin

After the build is complete, copy the `terraform-provider-solidfire` binary into the same path as your `terraform` binary, and re-run `terraform init`.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH` matches the architecture of your machine) file should contain a SHA256 sum that matches the local plugin. Run `shasum -a 256` on the binary to verify the values match.

## Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the [issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that work is not being duplicated. For further clarification, you can also ask in a new issue.

[gh-issues]: https://github.com/solidfire/terraform-provider-solidfire/issues
[gh-prs]: https://github.com/solidfire/terraform-provider-solidfire/pulls

See [Building the Provider](#building-the-provider) for details on building the provider.

## Testing the Provider

**NOTE:** Testing the SolidFire provider is currently a complex operation as it requires having an SolidFire endpoint to test against, which should be hosting a standard configuration for a HCI or SolidFire cluster. If you have a NetApp Support account, you may instead download Element Demo VM 12 from the Tools section, and deploy a singleton VM-based SolidFire cluster on native or nested VMware ESXi.

### Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment variables to run. See the individual `*_test.go` files in the [`elementsw/`](elementsw/) directory for more details. The next section also describes how you can manage a configuration file of the test environment variables.

#### Using the `.tf-elementsw-devrc.mk` file

The [`tf-elementsw-devrc.mk.example`](tf-elementsw-devrc.mk.example) file contains an up-to-date list of environment variables required to run the acceptance tests. Copy this to `$HOME/.tf-elementsw-devrc.mk` and change the permissions to something more secure (ie: `chmod 600 $HOME/.tf-elementsw-devrc.mk`), and configure the variables accordingly.

### Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the `TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccElementSwVolume"
```

This following example would run all of the acceptance tests matching `TestAccElementSwVolume`. Change this for the specific tests you want to run.

### Walk-through example

If you are not building from source or want to download provider from online Terraform Registry, please refer to README in the subdirectory examples/elementsw.

#### Installing Go and Terraform

```sh
bash
mkdir tf_na_elementsw
cd tf_na_elementsw

# if you want a private installation, use
export GO_INSTALL_DIR=`pwd`/go_install
mkdir $GO_INSTALL_DIR
# otherwise, go recommends to use
export GO_INSTALL_DIR=/usr/local

curl -O https://dl.google.com/go/go1.24.5.linux-amd64.tar.gz
tar -C $GO_INSTALL_DIR -xvf go1.24.5.linux-amd64.tar.gz

export PATH=$PATH:$GO_INSTALL_DIR/go/bin

curl -O https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip
unzip terraform_1.5.7_linux_amd64.zip
mv terraform $GO_INSTALL_DIR/go/bin

# you may want to add/update some of the exported variables in ~/.bashrc or other
```

#### Installing dependencies

```sh
# make sure git is installed
which git

export GOPATH=`pwd`
go get github.com/fatih/structs
go get github.com/hashicorp/terraform
go get github.com/sirupsen/logrus
go get github.com/x-cray/logrus-prefixed-formatter
```

Note getting the Terraform package also builds and installs Terraform in `$GOPATH/bin`.
The version in `go/bin` is a stable release.

#### Cloning the NetApp provider repository and building the provider

```sh
mkdir -p $GOPATH/src/github.com/scaleoutseasn
cd $GOPATH/src/github.com/scaleoutsean
git clone https://github.com/scaleoutsean/terraform-provider-solidfire.git
cd terraform-provider-solidfire
make build
mv $GOPATH/bin/terraform-provider-solidfire $GO_INSTALL_DIR/go/bin
```

The build step will install the provider in the `$GOPATH/bin` directory. Copy it to `/usr/share/terraform/providers/scaleoutsean.github.io/` and load it with:

```hcl
required_providers {
  netapp-elementsw = {
    version = "0.2.0"
    source = "scaleoutsean.github.io/netapp-elementsw"
  }
}
```

#### Sanity check

```sh
cd examples/elementsw/
terraform init
```

This should indicate `Terraform has been successfully initialized!`
