package elementsw

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceElementSwReplicationCluster manages SolidFire cluster pairing (replication)
func resourceElementSwReplicationCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwReplicationClusterCreate,
		Read:   resourceElementSwReplicationClusterRead,
		Delete: resourceElementSwReplicationClusterDelete,
		Schema: map[string]*schema.Schema{
			// Workflow 1: Manual/Key-based
			"pairing_key": {
				Type:        schema.TypeString,
				Optional:    true,
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

// resourceElementSwReplicationClusterCreate implements both pairing workflows
func resourceElementSwReplicationClusterCreate(d *schema.ResourceData, meta interface{}) error {
	// Always require target_cluster
	targetList, targetOk := d.GetOk("target_cluster")
	if !targetOk {
		return fmt.Errorf("target_cluster must be provided")
	}
	targetConn := expandClusterConnection(targetList)
	if targetConn == nil {
		return fmt.Errorf("invalid target_cluster connection info")
	}
	targetClient, err := NewSolidFireClient(targetConn.Endpoint, targetConn.Username, targetConn.Password)
	if err != nil {
		return fmt.Errorf("failed to create target cluster client: %w", err)
	}

	// Workflow 1: Manual/Key-based
	if v, ok := d.GetOk("pairing_key"); ok && v.(string) != "" {
		pairResp, err := targetClient.CompleteClusterPairing(v.(string))
		if err != nil {
			return fmt.Errorf("CompleteClusterPairing failed: %w", err)
		}
		d.SetId(fmt.Sprintf("%d", pairResp.ClusterPairID))
		_ = d.Set("cluster_pair_id", pairResp.ClusterPairID)
		return resourceElementSwReplicationClusterRead(d, meta)
	}

	// Workflow 2: Automated
	if sourceList, sourceOk := d.GetOk("source_cluster"); sourceOk {
		sourceConn := expandClusterConnection(sourceList)
		if sourceConn == nil {
			return fmt.Errorf("invalid source_cluster connection info")
		}
		sourceClient, err := NewSolidFireClient(sourceConn.Endpoint, sourceConn.Username, sourceConn.Password)
		if err != nil {
			return fmt.Errorf("failed to create source cluster client: %w", err)
		}
		keyResp, err := sourceClient.StartClusterPairing()
		if err != nil {
			return fmt.Errorf("StartClusterPairing failed: %w", err)
		}
		pairResp, err := targetClient.CompleteClusterPairing(keyResp.PairingKey)
		if err != nil {
			return fmt.Errorf("CompleteClusterPairing failed: %w", err)
		}
		d.SetId(fmt.Sprintf("%d", pairResp.ClusterPairID))
		_ = d.Set("cluster_pair_id", pairResp.ClusterPairID)
		return resourceElementSwReplicationClusterRead(d, meta)
	}

	return fmt.Errorf("you must provide either pairing_key or source_cluster info (target_cluster is always required)")
}

// SolidFireClient is a minimal JSON-RPC client for SolidFire API.
type SolidFireClient struct {
	endpoint   string
	username   string
	password   string
	httpClient *http.Client
}

// NewSolidFireClient returns a new API client with TLS 1.2+ enforced.
func NewSolidFireClient(endpoint, username, password string) (*SolidFireClient, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}
	return &SolidFireClient{
		endpoint:   endpoint,
		username:   username,
		password:   password,
		httpClient: &http.Client{Transport: tr, Timeout: 30 * time.Second},
	}, nil
}

// jsonRPCRequest defines the JSON-RPC request envelope.
type jsonRPCRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	ID     int         `json:"id"`
}

// jsonRPCResponse defines the JSON-RPC response envelope.
type jsonRPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// doRPC sends a JSON-RPC request and decodes the result.
func (c *SolidFireClient) doRPC(method string, params interface{}, result interface{}) error {
	reqBody, err := json.Marshal(jsonRPCRequest{Method: method, Params: params, ID: 1})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var rpcResp jsonRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return fmt.Errorf("invalid JSON-RPC response: %w", err)
	}
	if rpcResp.Error != nil {
		return fmt.Errorf("SolidFire API error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	if result != nil {
		if err := json.Unmarshal(rpcResp.Result, result); err != nil {
			return fmt.Errorf("failed to decode result: %w", err)
		}
	}
	return nil
}

// StartClusterPairingResponse holds the pairing key from StartClusterPairing API.
type StartClusterPairingResponse struct {
	PairingKey string `json:"pairingKey"`
}

// StartClusterPairing calls the API to start cluster pairing on the source cluster.
func (c *SolidFireClient) StartClusterPairing() (*StartClusterPairingResponse, error) {
	var resp StartClusterPairingResponse
	if err := c.doRPC("StartClusterPairing", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CompleteClusterPairingResponse holds the clusterPairID.
type CompleteClusterPairingResponse struct {
	ClusterPairID int `json:"clusterPairID"`
}

// CompleteClusterPairing calls the API to complete cluster pairing on the target cluster.
func (c *SolidFireClient) CompleteClusterPairing(pairingKey string) (*CompleteClusterPairingResponse, error) {
	var resp CompleteClusterPairingResponse
	params := map[string]interface{}{"pairingKey": pairingKey}
	if err := c.doRPC("CompleteClusterPairing", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// resourceElementSwReplicationClusterRead reads the current state of the cluster pairing.
func resourceElementSwReplicationClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	// List cluster pairs
	var listResp struct {
		ClusterPairs []struct {
			ClusterName     string `json:"clusterName"`
			ClusterPairID   int    `json:"clusterPairID"`
			Status          string `json:"status"`
		} `json:"clusterPairs"`
	}
	resp, err := client.CallAPIMethod("ListClusterPairs", map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to list cluster pairs: %w", err)
	}
	if err := json.Unmarshal(*resp, &listResp); err != nil {
		return fmt.Errorf("failed to unmarshal ListClusterPairs response: %w", err)
	}

	clusterPairID := d.Get("cluster_pair_id").(int)
	for _, pair := range listResp.ClusterPairs {
		if pair.ClusterPairID == clusterPairID {
			_ = d.Set("cluster_name", pair.ClusterName)
			_ = d.Set("status", pair.Status)
			return nil
		}
	}
	d.SetId("") // Cluster pair not found
	return nil
}

// resourceElementSwReplicationClusterDelete removes the cluster pairing.
func resourceElementSwReplicationClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	clusterPairID := d.Get("cluster_pair_id").(int)

	// Remove cluster pair
	params := map[string]interface{}{
		"clusterPairID": clusterPairID,
	}
	_, err := client.CallAPIMethod("RemoveClusterPair", params)
	if err != nil {
		return fmt.Errorf("failed to remove cluster pair: %w", err)
	}
	d.SetId("")
	return nil
}