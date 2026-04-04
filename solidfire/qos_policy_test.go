package solidfire

import (
	"os"
	"testing"

	"github.com/scaleoutsean/solidfire-go/sdk"
	"github.com/stretchr/testify/assert"
)

func getTestClient() *Client {
	c := &configStuct{
		User:            os.Getenv("SOLIDFIRE_USERNAME"),
		Password:        os.Getenv("SOLIDFIRE_PASSWORD"),
		ElementSwServer: os.Getenv("SOLIDFIRE_SERVER"),
		APIVersion:      os.Getenv("SOLIDFIRE_API_VERSION"),
	}
	client, err := c.clientFun()
	if err != nil {
		panic(err)
	}
	return client
}

func TestQoSPolicyLifecycle(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping TestQoSPolicyLifecycle; TF_ACC not set")
	}
	client := getTestClient()

	qos := sdk.QoS{
		MinIOPS:   100,
		MaxIOPS:   200,
		BurstIOPS: 300,
	}

	// Create
	id, err := client.CreateQoSPolicy("test-policy-lifecycle", qos)
	assert.NoError(t, err)
	assert.True(t, id > 0)

	// Get
	result, err := client.GetQoSPolicy(id)
	assert.NoError(t, err)
	if result != nil {
		assert.Equal(t, id, result.QosPolicyID)
		assert.Equal(t, "test-policy-lifecycle", result.Name)
	}

	// Modify
	qosUpdate := sdk.QoS{
		MinIOPS:   150,
		MaxIOPS:   250,
		BurstIOPS: 350,
	}
	err = client.ModifyQoSPolicy(id, "updated-policy-lifecycle", qosUpdate)
	assert.NoError(t, err)

	// List
	listResult, err := client.ListQoSPolicies()
	assert.NoError(t, err)
	assert.NotEmpty(t, listResult)

	// Delete
	err = client.DeleteQoSPolicy(id)
	assert.NoError(t, err)
}
