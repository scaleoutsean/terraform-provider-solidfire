package elementsw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

type initiator struct {
	Name                string      `json:"name"`
	Alias               string      `json:"alias"`
	Attributes          interface{} `json:"attributes"`
	VolumeAccessGroupID int64       `json:"volumeAccessGroupID"`
	InitiatorID         int64       `json:"initiatorID"`
}

func (c *Client) getInitiatorByID(id string) (initiator, error) {
	convID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return initiator{}, err
	}

	req := sdk.ListInitiatorsRequest{
		Initiators: []int64{convID},
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListInitiators(context.TODO(), &req)
	if sdkErr != nil {
		return initiator{}, sdkErr
	}

	if len(res.Initiators) != 1 {
		return initiator{}, fmt.Errorf("expected one initiator to be found. response contained %v results", len(res.Initiators))
	}

	var init initiator
	init.Name = res.Initiators[0].InitiatorName
	init.Alias = res.Initiators[0].Alias
	init.Attributes = res.Initiators[0].Attributes
	init.InitiatorID = res.Initiators[0].InitiatorID
	if len(res.Initiators[0].VolumeAccessGroups) > 0 {
		init.VolumeAccessGroupID = res.Initiators[0].VolumeAccessGroups[0]
	}

	return init, nil
}
