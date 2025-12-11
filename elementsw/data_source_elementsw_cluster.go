package elementsw

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceElementSwCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementSwClusterRead,
		Schema: map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Computed: true},
			"unique_id": {Type: schema.TypeString, Computed: true},
			"cluster_version": {Type: schema.TypeString, Computed: true},
			"cluster_api_version": {Type: schema.TypeString, Computed: true},
			"mvip": {Type: schema.TypeString, Computed: true},
			"svip": {Type: schema.TypeString, Computed: true},
		},
	}
}

type getClusterInfoResponse struct {
	Result struct {
		ClusterInfo struct {
			Name     string `json:"name"`
			UniqueID string `json:"uniqueID"`
			Mvip     string `json:"mvip"`
			Svip     string `json:"svip"`
		} `json:"clusterInfo"`
	} `json:"result"`
}

type getClusterVersionInfoResponse struct {
	Result struct {
		ClusterVersion     string `json:"clusterVersion"`
		ClusterAPIVersion  string `json:"clusterAPIVersion"`
	} `json:"result"`
}

func dataSourceElementSwClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	// GetClusterInfo
	infoRaw, err := client.CallAPIMethod("GetClusterInfo", map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("GetClusterInfo failed: %v", err)
	}
	var info getClusterInfoResponse
	if err := json.Unmarshal([]byte(*infoRaw), &info); err != nil {
		return fmt.Errorf("failed to decode GetClusterInfo: %v", err)
	}

	// GetClusterVersionInfo
	verRaw, err := client.CallAPIMethod("GetClusterVersionInfo", map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("GetClusterVersionInfo failed: %v", err)
	}
	var ver getClusterVersionInfoResponse
	if err := json.Unmarshal([]byte(*verRaw), &ver); err != nil {
		return fmt.Errorf("failed to decode GetClusterVersionInfo: %v", err)
	}

	d.SetId(info.Result.ClusterInfo.UniqueID)
	d.Set("name", info.Result.ClusterInfo.Name)
	d.Set("unique_id", info.Result.ClusterInfo.UniqueID)
	d.Set("mvip", info.Result.ClusterInfo.Mvip)
	d.Set("svip", info.Result.ClusterInfo.Svip)
	d.Set("cluster_version", ver.Result.ClusterVersion)
	d.Set("cluster_api_version", ver.Result.ClusterAPIVersion)

	return nil
}
