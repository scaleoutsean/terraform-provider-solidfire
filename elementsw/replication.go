package elementsw

import (
	"context"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

func (c *Client) StartClusterPairing() (*sdk.StartClusterPairingResult, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.StartClusterPairing(context.TODO())
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res, nil
}

func (c *Client) CompleteClusterPairing(key string) (*sdk.CompleteClusterPairingResult, error) {
	req := sdk.CompleteClusterPairingRequest{
		ClusterPairingKey: key,
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.CompleteClusterPairing(context.TODO(), &req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res, nil
}

func (c *Client) ListClusterPairs() ([]sdk.PairedCluster, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListClusterPairs(context.TODO())
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res.ClusterPairs, nil
}

func (c *Client) RemoveClusterPair(id int64) error {
	req := sdk.RemoveClusterPairRequest{
		ClusterPairID: id,
	}
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.RemoveClusterPair(context.TODO(), &req)
	return sdkErr
}

func (c *Client) StartVolumePairing(volumeID int64, mode string) (*sdk.StartVolumePairingResult, error) {
	req := sdk.StartVolumePairingRequest{
		VolumeID: volumeID,
		Mode:     mode,
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.StartVolumePairing(context.TODO(), &req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res, nil
}

func (c *Client) CompleteVolumePairing(volumeID int64, key string) error {
	req := sdk.CompleteVolumePairingRequest{
		VolumeID:         volumeID,
		VolumePairingKey: key,
	}
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.CompleteVolumePairing(context.TODO(), &req)
	return sdkErr
}

func (c *Client) ListActivePairedVolumes() ([]sdk.Volume, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListActivePairedVolumes(context.TODO(), &sdk.ListActivePairedVolumesRequest{})
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res.Volumes, nil
}

func (c *Client) ModifyVolumePair(req *sdk.ModifyVolumePairRequest) error {
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.ModifyVolumePair(context.TODO(), req)
	return sdkErr
}

func (c *Client) RemoveVolumePair(volumeID int64) error {
	req := sdk.RemoveVolumePairRequest{
		VolumeID: volumeID,
	}
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.RemoveVolumePair(context.TODO(), &req)
	return sdkErr
}
