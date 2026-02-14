package elementsw

import (
	"context"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

func (c *Client) CreateSnapshot(req *sdk.CreateSnapshotRequest) (*sdk.CreateSnapshotResult, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.CreateSnapshot(context.TODO(), req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res, nil
}

func (c *Client) DeleteSnapshot(id int64) error {
	req := sdk.DeleteSnapshotRequest{
		SnapshotID: id,
	}
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.DeleteSnapshot(context.TODO(), &req)
	return sdkErr
}

func (c *Client) ListSnapshots(volumeID int64) ([]sdk.Snapshot, error) {
	req := sdk.ListSnapshotsRequest{}
	if volumeID > 0 {
		req.VolumeID = volumeID
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListSnapshots(context.TODO(), &req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res.Snapshots, nil
}

func (c *Client) CreateGroupSnapshot(req *sdk.CreateGroupSnapshotRequest) (*sdk.CreateGroupSnapshotResult, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.CreateGroupSnapshot(context.TODO(), req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res, nil
}

func (c *Client) ModifySnapshot(req *sdk.ModifySnapshotRequest) error {
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.ModifySnapshot(context.TODO(), req)
	return sdkErr
}

func (c *Client) ModifyGroupSnapshot(req *sdk.ModifyGroupSnapshotRequest) error {
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.ModifyGroupSnapshot(context.TODO(), req)
	return sdkErr
}

func (c *Client) DeleteGroupSnapshot(id int64, saveMembers bool) error {
	req := sdk.DeleteGroupSnapshotRequest{
		GroupSnapshotID: id,
		SaveMembers:     saveMembers,
	}
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.DeleteGroupSnapshot(context.TODO(), &req)
	return sdkErr
}

func (c *Client) ListGroupSnapshots(volumeIDs []int64) ([]sdk.GroupSnapshot, error) {
	req := sdk.ListGroupSnapshotsRequest{}
	if len(volumeIDs) > 0 {
		req.Volumes = volumeIDs
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListGroupSnapshots(context.TODO(), &req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res.GroupSnapshots, nil
}
