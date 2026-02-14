package elementsw

import (
	"context"
	"fmt"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

func (c *Client) ListVolumes(volumeIDs []int64) ([]sdk.Volume, error) {
	req := sdk.ListVolumesRequest{}
	if len(volumeIDs) > 0 {
		req.VolumeIDs = volumeIDs
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListVolumes(context.TODO(), &req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res.Volumes, nil
}

func (c *Client) GetVolume(volumeID int64) (*sdk.Volume, error) {
	vols, err := c.ListVolumes([]int64{volumeID})
	if err != nil {
		return nil, err
	}
	if len(vols) == 0 {
		return nil, fmt.Errorf("volume %d not found", volumeID)
	}
	return &vols[0], nil
}

func (c *Client) ListVolumesForAccount(accountID int64) ([]sdk.Volume, error) {
	req := sdk.ListVolumesForAccountRequest{
		AccountID: accountID,
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListVolumesForAccount(context.TODO(), &req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res.Volumes, nil
}
