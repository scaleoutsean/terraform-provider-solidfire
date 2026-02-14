package elementsw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

type volumeAccessGroup struct {
	VolumeAccessGroupID int64    `json:"volumeAccessGroupID"`
	Name                string   `json:"name"`
	Initiators          []string `json:"initiators"`
	Volumes             []int64  `json:"volumes"`
	ID                  int64    `json:"id"`
}

func (c *Client) getVolumeAccessGroupByID(id string) (volumeAccessGroup, error) {
	convID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return volumeAccessGroup{}, err
	}

	req := sdk.ListVolumeAccessGroupsRequest{
		VolumeAccessGroups: []int64{convID},
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListVolumeAccessGroups(context.TODO(), &req)
	if sdkErr != nil {
		return volumeAccessGroup{}, sdkErr
	}

	if len(res.VolumeAccessGroups) != 1 {
		return volumeAccessGroup{}, fmt.Errorf("expected one volume access group to be found")
	}

	vag := res.VolumeAccessGroups[0]
	return volumeAccessGroup{
		VolumeAccessGroupID: vag.VolumeAccessGroupID,
		Name:                vag.Name,
		Initiators:          vag.Initiators,
		Volumes:             vag.Volumes,
		ID:                  vag.VolumeAccessGroupID,
	}, nil
}
