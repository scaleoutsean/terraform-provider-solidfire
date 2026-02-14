package elementsw

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceElementSwCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementSwClusterRead,
		Schema: map[string]*schema.Schema{
			"name":                {Type: schema.TypeString, Computed: true},
			"unique_id":           {Type: schema.TypeString, Computed: true},
			"cluster_version":     {Type: schema.TypeString, Computed: true},
			"cluster_api_version": {Type: schema.TypeString, Computed: true},
			"mvip":                {Type: schema.TypeString, Computed: true},
			"svip":                {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceElementSwClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	// GetClusterInfo
	info, err := client.GetClusterInfo()
	if err != nil {
		return fmt.Errorf("GetClusterInfo failed: %v", err)
	}

	// GetClusterVersionInfo
	ver, err := client.GetClusterVersionInfo()
	if err != nil {
		return fmt.Errorf("GetClusterVersionInfo failed: %v", err)
	}

	d.SetId(info.ClusterInfo.UniqueID)
	d.Set("name", info.ClusterInfo.Name)
	d.Set("unique_id", info.ClusterInfo.UniqueID)
	d.Set("mvip", info.ClusterInfo.Mvip)
	d.Set("svip", info.ClusterInfo.Svip)
	d.Set("cluster_version", ver.ClusterVersion)
	d.Set("cluster_api_version", ver.ClusterAPIVersion)

	return nil
}
