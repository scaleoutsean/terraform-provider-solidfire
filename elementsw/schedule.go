package elementsw

import (
	"context"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

func (c *Client) CreateSchedule(req *sdk.CreateScheduleRequest) (int64, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.CreateSchedule(context.TODO(), req)
	if sdkErr != nil {
		return 0, sdkErr
	}
	return res.ScheduleID, nil
}

func (c *Client) GetSchedule(id int64) (*sdk.Schedule, error) {
	// ListSchedules and filter by ID since GetSchedule return type varies
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListSchedules(context.TODO())
	if sdkErr != nil {
		return nil, sdkErr
	}
	for _, s := range res.Schedules {
		if s.ScheduleID == id {
			return &s, nil
		}
	}
	return nil, nil
}

func (c *Client) ModifySchedule(req *sdk.ModifyScheduleRequest) error {
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.ModifySchedule(context.TODO(), req)
	return sdkErr
}

func (c *Client) ListSchedules() ([]sdk.Schedule, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListSchedules(context.TODO())
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res.Schedules, nil
}
