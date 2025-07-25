package elementsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQoSPolicy(t *testing.T) {
	client := &Client{} // TODO: mock or initialize with test config
	request := getQoSPolicyRequest{QoSPolicyID: 1}

	result, err := client.GetQoSPolicy(request)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.QoSPolicy.QoSPolicyID)
	assert.NotEmpty(t, result.QoSPolicy.Name)
	assert.NotNil(t, result.QoSPolicy.QoS)
}

func TestListQoSPolicies(t *testing.T) {
	client := &Client{} // TODO: mock or initialize with test config
	request := listQoSPoliciesRequest{}

	result, err := client.ListQoSPolicies(request)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.QoSPolicies)
}
