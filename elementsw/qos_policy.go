package elementsw

import (
	"context"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

func (c *Client) CreateQoSPolicy(name string, qos sdk.QoS) (int64, error) {
	req := sdk.CreateQoSPolicyRequest{
		Name: name,
		Qos:  qos,
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.CreateQoSPolicy(context.TODO(), &req)
	if sdkErr != nil {
		return 0, sdkErr
	}
	return res.QosPolicy.QosPolicyID, nil
}

func (c *Client) GetQoSPolicy(id int64) (*sdk.QoSPolicy, error) {
	req := sdk.GetQoSPolicyRequest{
		QosPolicyID: id,
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.GetQoSPolicy(context.TODO(), &req)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return &res.QosPolicy, nil
}

func (c *Client) ModifyQoSPolicy(id int64, name string, qos sdk.QoS) error {
	req := sdk.ModifyQoSPolicyRequest{
		QosPolicyID: id,
		Name:        name,
		Qos:         qos,
	}
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.ModifyQoSPolicy(context.TODO(), &req)
	return sdkErr
}

func (c *Client) DeleteQoSPolicy(id int64) error {
	req := sdk.DeleteQoSPolicyRequest{
		QosPolicyID: id,
	}
	c.initOnce.Do(c.init)
	_, sdkErr := c.sdkClient.DeleteQoSPolicy(context.TODO(), &req)
	return sdkErr
}

func (c *Client) ListQoSPolicies() ([]sdk.QoSPolicy, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListQoSPolicies(context.TODO())
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res.QosPolicies, nil
}
