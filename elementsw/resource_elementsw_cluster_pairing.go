package elementsw

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

// resourceElementSwClusterPairing manages SolidFire cluster pairing (replication)
func resourceElementSwClusterPairing() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwClusterPairingCreate,
		Read:   resourceElementSwClusterPairingRead,
		Delete: resourceElementSwClusterPairingDelete,
		Schema: map[string]*schema.Schema{
			// Workflow 1: Manual/Key-based
			"pairing_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Pairing key generated from StartClusterPairing on the source cluster.",
			},
			// Workflow 2: Automated
			"source_cluster": clusterConnectionSchema("Source cluster for pairing (API endpoint, username, password)"),
			// Always required: target cluster
			"target_cluster": clusterConnectionSchemaRequired("Target cluster for pairing (API endpoint, username, password)"),
			// Common outputs
			"cluster_pair_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// clusterConnectionSchema returns an optional schema for cluster connection info
func clusterConnectionSchema(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"endpoint": {Type: schema.TypeString, Required: true},
				"username": {Type: schema.TypeString, Required: true},
				"password": {Type: schema.TypeString, Required: true, Sensitive: true},
			},
		},
	}
}

// clusterConnectionSchemaRequired returns a required schema for cluster connection info
func clusterConnectionSchemaRequired(desc string) *schema.Schema {
	s := clusterConnectionSchema(desc)
	s.Optional = false
	s.Required = true
	return s
}

// ClusterConnection holds endpoint/username/password for a cluster
type ClusterConnection struct {
	Endpoint string
	Username string
	Password string
}

// expandClusterConnection extracts endpoint/username/password from a schema.TypeList
func expandClusterConnection(list interface{}) *ClusterConnection {
	arr, ok := list.([]interface{})
	if !ok || len(arr) == 0 || arr[0] == nil {
		return nil
	}
	m, ok := arr[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return &ClusterConnection{
		Endpoint: m["endpoint"].(string),
		Username: m["username"].(string),
		Password: m["password"].(string),
	}
}

// createSFClientFromConn creates an SDK client from ClusterConnection
func createSFClientFromConn(conn *ClusterConnection) (*sdk.SFClient, error) {
	u, err := url.Parse(conn.Endpoint)
	if err != nil {
		return nil, err
	}
	// Extract version from path if matching /json-rpc/VERSION
	version := "12.5" // default
	parts := strings.Split(u.Path, "/")
	if len(parts) >= 3 && parts[1] == "json-rpc" {
		version = parts[2]
	}

	client := &sdk.SFClient{}
	// Note: using Host only as sdk.Connect builds URL
	client.Connect(context.TODO(), u.Host, version, conn.Username, conn.Password)
	return client, nil
}

// resourceElementSwClusterPairingCreate implements both pairing workflows
func resourceElementSwClusterPairingCreate(d *schema.ResourceData, meta interface{}) error {
	// Always require target_cluster
	targetList, targetOk := d.GetOk("target_cluster")
	if !targetOk {
		return fmt.Errorf("target_cluster must be provided")
	}
	targetConn := expandClusterConnection(targetList)
	if targetConn == nil {
		return fmt.Errorf("invalid target_cluster connection info")
	}
	targetClient, err := createSFClientFromConn(targetConn)
	if err != nil {
		return fmt.Errorf("failed to create target cluster client: %w", err)
	}

	// Workflow 1: Manual/Key-based
	if v, ok := d.GetOk("pairing_key"); ok && v.(string) != "" {
		req := sdk.CompleteClusterPairingRequest{
			ClusterPairingKey: v.(string),
		}
		pairResp, sdkErr := targetClient.CompleteClusterPairing(context.TODO(), &req)
		if sdkErr != nil {
			return fmt.Errorf("CompleteClusterPairing failed: %s", sdkErr.Detail)
		}
		d.SetId(fmt.Sprintf("%d", pairResp.ClusterPairID))
		_ = d.Set("cluster_pair_id", int(pairResp.ClusterPairID))
		return resourceElementSwClusterPairingRead(d, meta)
	}

	// Workflow 2: Automated
	if sourceList, sourceOk := d.GetOk("source_cluster"); sourceOk {
		sourceConn := expandClusterConnection(sourceList)
		if sourceConn == nil {
			return fmt.Errorf("invalid source_cluster connection info")
		}
		sourceClient, err := createSFClientFromConn(sourceConn)
		if err != nil {
			return fmt.Errorf("failed to create source cluster client: %w", err)
		}
		keyResp, sdkErr := sourceClient.StartClusterPairing(context.TODO())
		if sdkErr != nil {
			return fmt.Errorf("StartClusterPairing failed: %s", sdkErr.Detail)
		}
		req := sdk.CompleteClusterPairingRequest{
			ClusterPairingKey: keyResp.ClusterPairingKey,
		}
		pairResp, sdkErr := targetClient.CompleteClusterPairing(context.TODO(), &req)
		if sdkErr != nil {
			return fmt.Errorf("CompleteClusterPairing failed: %s", sdkErr.Detail)
		}
		d.SetId(fmt.Sprintf("%d", pairResp.ClusterPairID))
		_ = d.Set("cluster_pair_id", int(pairResp.ClusterPairID))
		return resourceElementSwClusterPairingRead(d, meta)
	}

	return fmt.Errorf("you must provide either pairing_key or source_cluster info (target_cluster is always required)")
}

// resourceElementSwClusterPairingRead reads the current state of the cluster pairing.
func resourceElementSwClusterPairingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	// List cluster pairs
	clusterPairs, err := client.ListClusterPairs()
	if err != nil {
		return fmt.Errorf("failed to list cluster pairs: %w", err)
	}

	clusterPairID := int64(d.Get("cluster_pair_id").(int))
	for _, pair := range clusterPairs {
		if pair.ClusterPairID == clusterPairID {
			_ = d.Set("cluster_name", pair.ClusterName)
			_ = d.Set("status", pair.Status)
			return nil
		}
	}
	d.SetId("") // Cluster pair not found
	return nil
}

// resourceElementSwClusterPairingDelete removes the cluster pairing.
func resourceElementSwClusterPairingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	clusterPairID := int64(d.Get("cluster_pair_id").(int))

	// Remove cluster pair
	err := client.RemoveClusterPair(clusterPairID)
	if err != nil {
		return fmt.Errorf("failed to remove cluster pair: %w", err)
	}
	d.SetId("")
	return nil
}
