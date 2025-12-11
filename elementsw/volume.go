package elementsw

import (
	"encoding/json"
	"log"

	"github.com/fatih/structs"
)


type listVolumesRequest struct {
	Volumes               []int `structs:"volumeIDs"`
	IncludeVirtualVolumes bool  `structs:"includeVirtualVolumes"`
}

type listVolumesResult struct {
	Volumes []volume `json:"volumes"`
}


// listVolumesByAccountIDRequest is used to list volumes for a given account and volume ID.
type listVolumesByAccountIDRequest struct {
	Accounts []int `structs:"accounts"`
}

type volume struct {
	Name     string `json:"name"`
	VolumeID int    `json:"volumeID"`
	Iqn      string `json:"iqn"`
}

func (c *Client) listVolumesByVolumeID(request listVolumesByAccountIDRequest) (listVolumesResult, error) {
	params := structs.Map(request)
	return c.getVolumesDetails(params)
}


func (c *Client) listVolumes(request listVolumesRequest) (listVolumesResult, error) {
	params := structs.Map(request)
	return c.getVolumesDetails(params)
}

func (c *Client) getVolumesDetails(params map[string]interface{}) (listVolumesResult, error) {

	response, err := c.CallAPIMethod("ListVolumes", params)
	if err != nil {
		log.Print("ListVolumes request failed")
		return listVolumesResult{}, err
	}

	var result listVolumesResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ListVolumes")
		return listVolumesResult{}, err
	}

	return result, nil
}

