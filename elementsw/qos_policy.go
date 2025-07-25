// ...existing code...

package elementsw

import (
	"encoding/json"
	"log"

	"fmt"

	"github.com/fatih/structs"
)

type createQoSPolicyRequest struct {
	Name string     `structs:"name"`
	QoS  qosDetails `structs:"qos"`
}

type createQoSPolicyResult struct {
	QoSPolicyID int `json:"qosPolicyID"`
}

func (c *Client) CreateQoSPolicy(request createQoSPolicyRequest) (createQoSPolicyResult, error) {
	// Validation
	if len(request.Name) < 1 || len(request.Name) > 40 {
		return createQoSPolicyResult{}, fmt.Errorf("name must be 1-40 alphanumeric characters")
	}
	// Optionally, add regex for alphanumeric check
	if request.QoS.MinIOPS < 50 || request.QoS.MinIOPS > 15000 {
		return createQoSPolicyResult{}, fmt.Errorf("minIOPS must be 50-15000")
	}
	if request.QoS.MaxIOPS < 51 || request.QoS.MaxIOPS > 50000 || request.QoS.MaxIOPS <= request.QoS.MinIOPS {
		return createQoSPolicyResult{}, fmt.Errorf("maxIOPS must be 51-50000 and greater than minIOPS")
	}
	if request.QoS.BurstIOPS < 100 || request.QoS.BurstIOPS > 20000 || request.QoS.BurstIOPS <= request.QoS.MaxIOPS {
		return createQoSPolicyResult{}, fmt.Errorf("burstIOPS must be 100-20000 and greater than maxIOPS")
	}

	params := structs.Map(request)
	response, err := c.CallAPIMethod("CreateQoSPolicy", params)
	if err != nil {
		log.Print("CreateQoSPolicy request failed")
		return createQoSPolicyResult{}, err
	}

	var result createQoSPolicyResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from CreateQoSPolicy")
		return createQoSPolicyResult{}, err
	}

	return result, nil
}

type getQoSPolicyRequest struct {
	QoSPolicyID int `structs:"qosPolicyID"`
}

type getQoSPolicyResult struct {
	QoSPolicy qosPolicy `json:"qosPolicy"`
}

func (c *Client) GetQoSPolicy(request getQoSPolicyRequest) (getQoSPolicyResult, error) {
	params := structs.Map(request)
	response, err := c.CallAPIMethod("GetQoSPolicy", params)
	if err != nil {
		log.Print("GetQoSPolicy request failed")
		return getQoSPolicyResult{}, err
	}

	var result getQoSPolicyResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from GetQoSPolicy")
		return getQoSPolicyResult{}, err
	}

	return result, nil
}

type listQoSPoliciesRequest struct {
	// No parameters for ListQoSPolicies
}

type qosPolicy struct {
	Name        string     `json:"name"`
	QoS         qosDetails `json:"qos"`
	QoSPolicyID int        `json:"qosPolicyID"`
	VolumeIDs   []int      `json:"volumeIDs"`
}

type listQoSPoliciesResult struct {
	QoSPolicies []qosPolicy `json:"qosPolicies"`
}

func (c *Client) ListQoSPolicies(request listQoSPoliciesRequest) (listQoSPoliciesResult, error) {
	params := structs.Map(request)
	response, err := c.CallAPIMethod("ListQoSPolicies", params)
	if err != nil {
		log.Print("ListQoSPolicies request failed")
		return listQoSPoliciesResult{}, err
	}

	var result listQoSPoliciesResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ListQoSPolicies")
		return listQoSPoliciesResult{}, err
	}

	return result, nil
}

type qosCurve map[string]int

type qosDetails struct {
	BurstIOPS int      `json:"burstIOPS"`
	BurstTime int      `json:"burstTime"`
	Curve     qosCurve `json:"curve"`
	MaxIOPS   int      `json:"maxIOPS"`
	MinIOPS   int      `json:"minIOPS"`
}
