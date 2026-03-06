# NetApp ElementSW v20.11 Example

Examples in resources.tf.example1 and resources.tf.example2 are designed to demonstrate the capabilities of the [Terraform
NetApp ElementSW Provider][ref-tf-elementsw].

[ref-tf-elementsw]: https://registry.terraform.io/providers/NetApp/netapp-elementsw/latest

## Requirements

* NetApp HCI, SolidFire or eSDS storage cluster (including Element Demo VM)
* Terraform client

## Getting Started

Clone the Git repository and change directory to the ElementSW examples directory:

```sh
git clone https://github.com/NetApp/terraform-provider-netapp-elementsw
cd terraform-provider-netapp-elementsw/examples/elementsw
```

**NOTE:** Before you continue make sure that volume names, sizes, IQN and other variables from the examples do not conflict with your production environment. Pay special attention when deleting resources because there is no undo. As mentioned in the main README file, you may download Element (SolidFire) Demo VM for safe experimenting.

### Example one: create an account and volumes for CHAP access

To try the first example, `resources.tf.example1`, copy the file to `resources.tf` and examine its contents including tenant and volume names so that you can adjust them if they conflict with your current environment.

Run `terraform init` to doownload NetApp ElementSW Provider.

On `terraform apply`, this example will perform the following:

* Set up a tenant account. This uses the `elementsw_account` resource.
* Creates two volumes for the account using the `elementsw_volume` resource.

`terraform apply` requires certain inputs. You can provide them in `terraform.tfvars` (see `terraform.tfvars.example`) or pass them from the CLI like so:

```sh
terraform apply \
  -var="elementsw_username=admin" \
  -var="elementsw_password=admin" \
  -var="elementsw_cluster=192.168.1.34"
```

On `terraform destroy`, all the resources will be deleted (volumes are purged, not just deleted) without the option to undo. You may need to provide the same variables as above - SolidFire cluster username, password and Management Virtual IP.

After first successful apply, make changes to `resources.tf` and run apply again.

If you want to try the second example, remember to destroy resources with `terraform destroy` and then copy the second example, `resources.tf.example2`, over `resources.tf` (if you had it from the first example). Without a clean-up you may encounter errors due to overlap between resources in the two examples.

### Example two: create an account and volumes VAG access

The second example demonstrates the use of Volume Access Group (VAG) and Initiator resource to creates two additional resources:

* Volume Access Group (VAG) for the volumes, using the `elementsw_volume_access_group` resource.
* Initiator tied to the VAG and the volumes using the `elementsw_initiator` resource.

It also does two things differently from the first example:

* Uses a list of volumes, which is simple but less flexible.
* Lets the SolidFire API to automatically generate tenant secrets - also simple, but less flexible.

Because some variables in this example have values set in `resources.tf` and some have defaults defined in `variables.tf`, the number of variables we have to provide via command line can be less than total number of required variables. For example, `elementsw_username` is already defined in `variables.tf` and `elementsw_initiator` in `resources.tf`, but we can still override the value of former through the CLI.

Like in the first example, check the values of variables in these files and change them to avoid any conflict with existing resources.

```sh
terraform apply \
  -var="elementsw_username=admin" \
  -var="elementsw_password=admin" \
  -var="elementsw_cluster=192.168.1.34" \
  -var="volume_name=testVol" \
  -var="volume_size_list=[1073742000,1073742000]"
```

Note that in this example `volume_size_list` defaults to `[]` (empty list) in order to avoid potential problems. You can change the default value if you want to change this behavior.

To destroy resources just created, run `terraform destroy` (you may need to provide the first three variables).

Descend to examples/elementsw subdirectory and use the sample file with variables to create `terraform.tfvars` and then edit the new file to match your environment:

```sh
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars
```

Now run `terraform plan` followed by `terraform apply`. You can omit most variables, but beware of security implications of having `elementsw_password` in plain text file. You may still choose to override certain default variables or variables set in `terraform.tfvars`, especially if they are similar or identical to existing resources.

Destroy with `terraform destroy`, the same way as before.

#### Overriding variable values from the CLI

This example only shows how values for two maps (QoS and IQN) can be provided from the CLI (Bash shell on Linux). Variations of this approach may be required for different OS.

```sh
terraform apply \
  -var="elementsw_username=admin" \
  -var="elementsw_password=admin" \
  -var="elementsw_cluster=192.168.1.34" \
  -var="volume_name=testVol" \
  -var="volume_size_list=[1073742000,1073742000,1073742000]" \
  -var="sectorsize_512e=false" \
  -var="qos={min=100,max=200,burst=300}" \
  -var="volume_name=dc1-testVol-master" \
  -var="elementsw_initiator={name=\"iqn.1998-01.com.vmware:test-cluster-000001\",alias=\"testNode1\"}" \
  -var="volume_group_name=testTenant" \
  -var="elementsw_tenant_name=testCluster01"
```

### Example three: Cluster Stats Data Source

This example demonstrates how to use the `elementsw_cluster_stats` data source to retrieve cluster capacity and performance metrics.

```hcl
data "elementsw_cluster_stats" "current" {}

output "cluster_health" {
  value = {
    volumes_total    = data.elementsw_cluster_stats.current.volume_count
    nodes_total      = data.elementsw_cluster_stats.current.node_count
    # Note: density is an average across all nodes. Individual nodes may have higher counts.
    density          = data.elementsw_cluster_stats.current.volumes_per_node
    efficiency_ratio = data.elementsw_cluster_stats.current.compression_factor
    current_iops     = data.elementsw_cluster_stats.current.metrics[0].actual_iops
    used_space       = data.elementsw_cluster_stats.current.capacity[0].used_space
  }
}
```

### Example four: Replication

This example demonstrates how to set up cluster and volume replication.

```hcl
# 1. Pair Clusters
resource "elementsw_replication_cluster" "dr_site" {
  source_cluster {
    endpoint = "https://10.10.10.10/json-rpc/12.5"
    username = "admin"
    password = "password"
  }
  target_cluster {
    endpoint = "https://10.20.20.20/json-rpc/12.5"
    username = "admin"
    password = "password"
  }
}

# 2. Create Volume on Source
resource "elementsw_volume" "primary" {
  name       = "primary-vol"
  account_id = 1
  total_size = 1000000000
  enable512e = true
}

# 3. Enable Replication for Volume
resource "elementsw_replication_volume" "dr_vol" {
  volume_id = elementsw_volume.primary.id
  
  # Optional: Provide target credentials to automatically complete pairing on the target side.
  # If omitted, pairing must be completed manually or via a separate process.
  target_cluster {
    endpoint = "https://10.20.20.20/json-rpc/12.5"
    username = "admin"
    password = "password"
  }
}
```

**Note:** Currently, if cluster pairing fails or is destroyed, there is no automatic cleanup of the "dangling" relationship on the remote cluster. You may need to manually remove the cluster pair on the other side or dangling relationship on the source. For example, the remote cluster could already have the maximum number of cluster relationships and depite correct functioning of this provider, cluster pairing with that cluster would fail.

How to use the Provider for site or cluster failover:

- `access` (the volume status) is a property of the volume itself that determines if it is the replication source (`readWrite`) or the target (`replicationTarget`).
- `mode` (`Async`, `Sync`, `SnapMirror`) and `paused` (boolean) are configurable attributes.

Pick a mode to set up replication, and simply swap the value of `access` properies of paired volumes to reverse the direction (A <- B).

Users of solidfire-csi, which uses volume IDs as volume handles and has account ID (tenant ID) storage classes, can easily set up replication and orchestrate site failover. Monitoring of cluster and volume pairings, replication delays and more is available in [SFC](https://github.com/scaleoutsean/sfc/).

```hcl
# Look up a volume created by K8s CSI (volumeHandle = ID)
data "elementsw_volume" "csi_vol" {
  volume_id = 123 
}

# Use it in your pairing resource
resource "elementsw_volume_pairing" "k8s_dr" {
  volume_id      = data.elementsw_volume.csi_vol.volume_id
  target_cluster = { ... }
}
```

`elementsw_volume` can be looked up by Name or ID without managing it in the current Terraform state.

```hcl
data "elementsw_volume" "from_k8s" {
  volume_id = 42
}

output "iqn" {
  value = data.elementsw_volume.from_k8s.iqn
}
```

### Add own validation rules

To implement own naming rules or conventions, feel free to create Terraform validation rules.

In this example we want to ensure that volume names begin with `dc1`.

```hcl
variable "volume_name" {
  type        = string
  description = "The Element volume name."

  validation {
    condition     = length(var.volume_name) > 2 && substr(var.volume_name, 0, 3) == "dc1"
    error_message = "The volume name string must begin with \"dc1\" and have 3 or more characters."
  }
}
```

`variables.tf` contains few other example of validation rules (acceptable volume sizes (min 1Gi, max 16TiB), initiator secrets, and volume QoS values).

### Extend

If you wish to extend the scope of this provider with minor features, Terraform [generic provisioners](https://www.terraform.io/docs/language/resources/provisioners/file.html) or vendor provisioners may be a convenient way to achieve that without developing in Go.

## Tests

Set the following variables:
- ELEMENTSW_USERNAME
- ELEMENTSW_PASSWORD
- ELEMENTSW_SERVER
- ELEMENTSW_API_VERSION

To run replication tests, also set:
- ELEMENTSW_SERVER_DR (Target cluster endpoint)
- ELEMENTSW_USERNAME_DR (Optional, defaults to ELEMENTSW_USERNAME)
- ELEMENTSW_PASSWORD_DR (Optional, defaults to ELEMENTSW_PASSWORD)

```sh
go test -v ./elementsw -run TestAccElementsw_FullCycle
```