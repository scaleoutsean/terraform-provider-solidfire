package elementsw

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceElementSwReplicationVolume manages SolidFire volume pairing (replication)
func resourceElementSwReplicationVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwReplicationVolumeCreate,
		Read:   resourceElementSwReplicationVolumeRead,
		Update: resourceElementSwReplicationVolumeUpdate,
		Delete: resourceElementSwReplicationVolumeDelete,
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the volume to be paired.",
			},
			"pairing_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The pairing key used to complete volume pairing.",
			},
			// Automated pairing support
			"target_cluster": clusterConnectionSchema("Target cluster for pairing (API endpoint, username, password)"),
		},
	}
}

func resourceElementSwReplicationVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	volumeID := d.Get("volume_id").(int)

	// 1. Start volume pairing on source
	params := map[string]interface{}{
		"volumeID": volumeID,
		"mode":     "Async",
	}
	resp, err := client.CallAPIMethod("StartVolumePairing", params)
	if err != nil {
		return fmt.Errorf("failed to start volume pairing: %w", err)
	}

	var pairingResp struct {
		VolumePairingKey string `json:"volumePairingKey"`
	}
	if err := json.Unmarshal(*resp, &pairingResp); err != nil {
		return fmt.Errorf("failed to unmarshal StartVolumePairing response: %w", err)
	}

	d.Set("pairing_key", pairingResp.VolumePairingKey)
	d.SetId(fmt.Sprintf("%d", volumeID))

	// 2. If target_cluster is provided, complete pairing on target
	if targetList, ok := d.GetOk("target_cluster"); ok {
		targetConn := expandClusterConnection(targetList)
		if targetConn != nil {
			// Create target client
			targetClient, err := NewSolidFireClient(targetConn.Endpoint, targetConn.Username, targetConn.Password)
			if err != nil {
				return fmt.Errorf("failed to create target cluster client: %w", err)
			}

			// We need the target volume ID.
			// Strategy: Get source volume name, find volume with same name on target.
			
			// Get source volume details
			var getVolResp struct {
				Volume struct {
					Name string `json:"name"`
				} `json:"volume"`
			}
			volParams := map[string]interface{}{"volumeID": volumeID}
			volRespBytes, err := client.CallAPIMethod("GetVolume", volParams)
			if err != nil {
				return fmt.Errorf("failed to get source volume details: %w", err)
			}
			if err := json.Unmarshal(*volRespBytes, &getVolResp); err != nil {
				return fmt.Errorf("failed to unmarshal GetVolume response: %w", err)
			}
			sourceVolName := getVolResp.Volume.Name

			// Find volume on target
			var listResp struct {
				Volumes []struct {
					VolumeID int    `json:"volumeID"`
					Name     string `json:"name"`
				} `json:"volumes"`
			}
			
			targetVolumeID := 0
			startID := 0
			for {
				listParams := map[string]interface{}{
					"startVolumeID": startID,
					"limit":         1000,
				}
				if err := targetClient.doRPC("ListActiveVolumes", listParams, &listResp); err != nil {
					return fmt.Errorf("failed to list volumes on target: %w", err)
				}
				
				if len(listResp.Volumes) == 0 {
					break
				}
				
				for _, v := range listResp.Volumes {
					if v.Name == sourceVolName {
						targetVolumeID = v.VolumeID
						break
					}
					if v.VolumeID > startID {
						startID = v.VolumeID
					}
				}
				if targetVolumeID != 0 {
					break
				}
				// If we got fewer than limit, we are done
				if len(listResp.Volumes) < 1000 {
					break
				}
				startID++ // Next batch
			}

			if targetVolumeID == 0 {
				return fmt.Errorf("target volume with name '%s' not found on target cluster", sourceVolName)
			}

			// Complete pairing on target
			completeParams := map[string]interface{}{
				"volumePairingKey": pairingResp.VolumePairingKey,
				"volumeID":         targetVolumeID,
			}
			if err := targetClient.doRPC("CompleteVolumePairing", completeParams, nil); err != nil {
				return fmt.Errorf("failed to complete volume pairing on target: %w", err)
			}
		}
	}

	return resourceElementSwReplicationVolumeRead(d, meta)
}

func resourceElementSwReplicationVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	volumeID := d.Get("volume_id").(int)

	// List active paired volumes
	params := map[string]interface{}{}
	resp, err := client.CallAPIMethod("ListActivePairedVolumes", params)
	if err != nil {
		return fmt.Errorf("failed to list active paired volumes: %w", err)
	}

	var listResp struct {
		Volumes []struct {
			VolumeID int `json:"volumeID"`
		} `json:"volumes"`
	}
	if err := json.Unmarshal(*resp, &listResp); err != nil {
		return fmt.Errorf("failed to unmarshal ListActivePairedVolumes response: %w", err)
	}

	for _, vol := range listResp.Volumes {
		if vol.VolumeID == volumeID {
			return nil // Volume is still paired
		}
	}

	d.SetId("") // Volume is no longer paired
	return nil
}

func resourceElementSwReplicationVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	volumeID := d.Get("volume_id").(int)

	// Modify volume pair
	params := map[string]interface{}{
		"volumeID":     volumeID,
		"pausedManual": false,
		"mode":         "Async",
	}
	_, err := client.CallAPIMethod("ModifyVolumePair", params)
	if err != nil {
		return fmt.Errorf("failed to modify volume pair: %w", err)
	}
	return resourceElementSwReplicationVolumeRead(d, meta)
}

func resourceElementSwReplicationVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	volumeID := d.Get("volume_id").(int)

	// Remove volume pair
	params := map[string]interface{}{
		"volumeID": volumeID,
	}
	_, err := client.CallAPIMethod("RemoveVolumePair", params)
	if err != nil {
		return fmt.Errorf("failed to remove volume pair: %w", err)
	}
	d.SetId("")
	return nil
}
