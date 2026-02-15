package elementsw

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

// resourceElementSwVolumePairing manages SolidFire volume pairing (replication)
func resourceElementSwVolumePairing() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwVolumePairingCreate,
		Read:   resourceElementSwVolumePairingRead,
		Update: resourceElementSwVolumePairingUpdate,
		Delete: resourceElementSwVolumePairingDelete,
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
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Async",
				Description: "The replication mode (Async, Sync, or SnapMirror).",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					valid := map[string]bool{
						"Async":      true,
						"Sync":       true,
						"SnapMirror": true,
					}
					if !valid[value] {
						errors = append(errors, fmt.Errorf("%q is not a valid replication mode", value))
					}
					return
				},
			},
			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to pause the volume pairing.",
			},
			// Automated pairing support
			"target_cluster": clusterConnectionSchema("Target cluster for pairing (API endpoint, username, password)"),
		},
	}
}

func resourceElementSwVolumePairingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	volumeID := int64(d.Get("volume_id").(int))
	mode := d.Get("mode").(string)

	// 1. Start volume pairing on source
	resp, err := client.StartVolumePairing(volumeID, mode)
	if err != nil {
		return fmt.Errorf("failed to start volume pairing: %w", err)
	}

	d.Set("pairing_key", resp.VolumePairingKey)
	d.SetId(fmt.Sprintf("%d", volumeID))

	// 2. If target_cluster is provided, complete pairing on target
	if targetList, ok := d.GetOk("target_cluster"); ok {
		targetConn := expandClusterConnection(targetList)
		if targetConn != nil {
			// Create target client
			targetSF, err := createSFClientFromConn(targetConn)
			if err != nil {
				return fmt.Errorf("failed to create target cluster client: %w", err)
			}
			targetClient := &Client{sdkClient: targetSF}
			// Important: Mark as initialized so targetClient.init() doesn't overwrite our SDK client
			targetClient.initOnce.Do(func() {})

			// We need the target volume ID.
			// Strategy: Get source volume name, find volume with same name on target.

			// Get source volume details
			vol, err := client.GetVolume(volumeID)
			if err != nil {
				return fmt.Errorf("failed to get source volume details: %w", err)
			}
			sourceVolName := vol.Name

			// Find volume on target with retry
			targetVolumeID := int64(0)
			for attempt := 0; attempt < 10; attempt++ {
				startID := int64(0)
				for {
					req := sdk.ListActiveVolumesRequest{
						StartVolumeID: startID,
						Limit:         1000,
					}
					volumes, err := targetClient.ListActiveVolumes(&req)
					if err != nil {
						return fmt.Errorf("failed to list volumes on target: %w", err)
					}

					if len(volumes) == 0 {
						break
					}

					for _, v := range volumes {
						if v.Name == sourceVolName {
							targetVolumeID = v.VolumeID
							break
						}
						if v.VolumeID > startID {
							startID = v.VolumeID
						}
					}
					if targetVolumeID != 0 || len(volumes) < 1000 {
						break
					}
					startID++
				}
				if targetVolumeID != 0 {
					break
				}
				time.Sleep(2 * time.Second)
			}

			if targetVolumeID == 0 {
				return fmt.Errorf("target volume with name '%s' not found on target cluster after retries", sourceVolName)
			}

			// Ensure target volume is set to replicationTarget mode before pairing
			ourlog.Infof("Setting target volume %d to 'replicationTarget' mode", targetVolumeID)
			modifyReq := &sdk.ModifyVolumeRequest{
				VolumeID: targetVolumeID,
				Access:   "replicationTarget",
			}
			err = targetClient.ModifyVolume(modifyReq)
			if err != nil {
				return fmt.Errorf("failed to set target volume to replicationTarget: %w", err)
			}

			// Complete pairing on target
			// Retry as cluster pairing might not be fully established (Connected) yet
			var lastErr error
			var success bool
			for i := 0; i < 20; i++ {
				err = targetClient.CompleteVolumePairing(targetVolumeID, resp.VolumePairingKey)
				if err == nil {
					success = true
					break
				}
				lastErr = err
				// xMVIPNotPaired means clusters are not yet paired or transitioning
				if !strings.Contains(err.Error(), "xMVIPNotPaired") {
					return fmt.Errorf("failed to complete volume pairing on target: %w", err)
				}
				ourlog.Infof("Waiting for cluster pairing to be ready (attempt %d/20)...", i+1)
				time.Sleep(3 * time.Second)
			}
			if !success {
				return fmt.Errorf("failed to complete volume pairing on target after retries: %w", lastErr)
			}
		}
	}

	return resourceElementSwVolumePairingRead(d, meta)
}

func resourceElementSwVolumePairingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	volumeID := int64(d.Get("volume_id").(int))

	// List active paired volumes
	vols, err := client.ListActivePairedVolumes()
	if err != nil {
		return fmt.Errorf("failed to list active paired volumes: %w", err)
	}

	for _, vol := range vols {
		if vol.VolumeID == volumeID {
			return nil // Volume is still paired
		}
	}

	d.SetId("") // Volume is no longer paired
	return nil
}

func resourceElementSwVolumePairingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	volumeID := int64(d.Get("volume_id").(int))

	// Modify volume pair
	req := sdk.ModifyVolumePairRequest{
		VolumeID:     volumeID,
		PausedManual: d.Get("paused").(bool),
		Mode:         d.Get("mode").(string),
	}
	err := client.ModifyVolumePair(&req)
	if err != nil {
		return fmt.Errorf("failed to modify volume pair: %w", err)
	}
	return resourceElementSwVolumePairingRead(d, meta)
}

func resourceElementSwVolumePairingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	volumeID := int64(d.Get("volume_id").(int))

	// Remove volume pair
	err := client.RemoveVolumePair(volumeID)
	if err != nil {
		return fmt.Errorf("failed to remove volume pair: %w", err)
	}
	d.SetId("")
	return nil
}
