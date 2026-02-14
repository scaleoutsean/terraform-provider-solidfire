package elementsw

import (
	"os"
	"testing"

	"github.com/scaleoutsean/solidfire-go/sdk"
	"github.com/stretchr/testify/assert"
)

func TestGetQoSPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestGetQoSPolicy; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config

	result, err := client.GetQoSPolicy(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.QosPolicyID)
	assert.NotEmpty(t, result.Name)
	assert.NotNil(t, result.Qos)
}

func TestListQoSPolicies(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestListQoSPolicies; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config

	result, err := client.ListQoSPolicies()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestCreateQoSPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestCreateQoSPolicy; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config
	qos := sdk.QoS{
		MinIOPS:   100,
		MaxIOPS:   200,
		BurstIOPS: 300,
	}
	id, err := client.CreateQoSPolicy("test-policy", qos)
	assert.NoError(t, err)
	assert.True(t, id > 0)
}

func TestModifyQoSPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestModifyQoSPolicy; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config
	qos := sdk.QoS{
		MinIOPS: 150,
	}
	err := client.ModifyQoSPolicy(1, "updated-policy", qos)
	assert.NoError(t, err)
}

func TestDeleteQoSPolicy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestDeleteQoSPolicy; TF_ACC not set")
	}
	client := &Client{} // TODO: mock or initialize with test config
	err := client.DeleteQoSPolicy(1)
	assert.NoError(t, err)
}
