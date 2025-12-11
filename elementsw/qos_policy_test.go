package elementsw

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQoSPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestGetQoSPolicy; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config
	request := getQoSPolicyRequest{QoSPolicyID: 1}

	result, err := client.GetQoSPolicy(request)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.QoSPolicy.QoSPolicyID)
	assert.NotEmpty(t, result.QoSPolicy.Name)
	assert.NotNil(t, result.QoSPolicy.QoS)
}

func TestListQoSPolicies(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestListQoSPolicies; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config
	request := listQoSPoliciesRequest{}

	result, err := client.ListQoSPolicies(request)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.QoSPolicies)
}

func TestCreateQoSPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestCreateQoSPolicy; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config
	request := createQoSPolicyRequest{
		Name: "test-policy",
		QoS: qosDetails{
			MinIOPS:   100,
			MaxIOPS:   200,
			BurstIOPS: 300,
			BurstTime: 60,
			Curve:     map[string]int{},
		},
	}
	result, err := client.CreateQoSPolicy(request)
	assert.NoError(t, err)
	assert.True(t, result.QoSPolicyID > 0)
}

func TestModifyQoSPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestModifyQoSPolicy; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config
	request := modifyQoSPolicyRequest{
		QoSPolicyID: 1,
		Name:       "updated-policy",
		QoS: qosDetails{
			MinIOPS:   150,
			MaxIOPS:   250,
			BurstIOPS: 350,
			BurstTime: 120,
			Curve:     map[string]int{},
		},
	}
	result, err := client.ModifyQoSPolicy(request)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.QoSPolicy.QoSPolicyID)
	assert.Equal(t, "updated-policy", result.QoSPolicy.Name)
}

func TestDeleteQoSPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestDeleteQoSPolicy; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config
	request := deleteQoSPolicyRequest{QoSPolicyID: 1}
	result, err := client.DeleteQoSPolicy(request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
