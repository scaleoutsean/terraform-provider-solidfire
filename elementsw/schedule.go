package elementsw

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

// ListSchedulesRequest is used to list all snapshot schedules
// No parameters required
// Corresponds to SolidFire Element API ListSchedules

type ListSchedulesRequest struct{}

type Schedule struct {
	Attributes           map[string]interface{} `json:"attributes"`
	HasError             bool                   `json:"hasError"`
	Hours                int                    `json:"hours"`
	LastRunStatus        string                 `json:"lastRunStatus"`
	LastRunTimeStarted   string                 `json:"lastRunTimeStarted"`
	Minutes              int                    `json:"minutes"`
	Monthdays            []int                  `json:"monthdays"`
	Paused               bool                   `json:"paused"`
	Recurring            bool                   `json:"recurring"`
	RunNextInterval      bool                   `json:"runNextInterval"`
	ScheduleID           int                    `json:"scheduleID"`
	ScheduleInfo         map[string]interface{} `json:"scheduleInfo"`
	ScheduleName         string                 `json:"scheduleName"`
	ScheduleType         string                 `json:"scheduleType"`
	StartingDate         *string                `json:"startingDate"`
	ToBeDeleted          bool                   `json:"toBeDeleted"`
	Weekdays             []int                  `json:"weekdays"`
}

type ListSchedulesResult struct {
	Schedules []Schedule `json:"schedules"`
}

func (c *Client) ListSchedules(request ListSchedulesRequest) (ListSchedulesResult, error) {
	params := map[string]interface{}{} // No params for ListSchedules
	response, err := c.CallAPIMethod("ListSchedules", params)
	if err != nil {
		log.Print("ListSchedules request failed")
		return ListSchedulesResult{}, err
	}

	var result ListSchedulesResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ListSchedules")
		return ListSchedulesResult{}, err
	}

	return result, nil
}

// GetScheduleRequest is used to get a specific snapshot schedule by ID
// Corresponds to SolidFire Element API GetSchedule

type GetScheduleRequest struct {
	ScheduleID int `structs:"scheduleID"`
}

type GetScheduleResult struct {
	Schedule Schedule `json:"schedule"`
}

func (c *Client) GetSchedule(request GetScheduleRequest) (GetScheduleResult, error) {
	params := map[string]interface{}{
		"scheduleID": request.ScheduleID,
	}
	response, err := c.CallAPIMethod("GetSchedule", params)
	if err != nil {
		log.Print("GetSchedule request failed")
		return GetScheduleResult{}, err
	}

	var result GetScheduleResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from GetSchedule")
		return GetScheduleResult{}, err
	}

	return result, nil
}

// CreateScheduleRequest is used to create a snapshot schedule
// Corresponds to SolidFire Element API CreateSchedule

type CreateScheduleRequest struct {
	Attributes      map[string]interface{} `structs:"attributes"`
	Hours           *int                  `structs:"hours,omitempty"`
	Minutes         *int                  `structs:"minutes,omitempty"`
	Paused          *bool                 `structs:"paused,omitempty"`
	Recurring       *bool                 `structs:"recurring,omitempty"`
	RunNextInterval *bool                 `structs:"runNextInterval,omitempty"`
	ScheduleName    string                `structs:"scheduleName"`
	ScheduleType    string                `structs:"scheduleType"`
	ScheduleInfo    map[string]interface{} `structs:"scheduleInfo"`
	StartingDate    *string               `structs:"startingDate,omitempty"`
	Monthdays       []int                 `structs:"monthdays,omitempty"`
	Weekdays        []map[string]int      `structs:"weekdays,omitempty"`
}

type CreateScheduleResult struct {
	Schedule Schedule `json:"schedule"`
}

func (c *Client) CreateSchedule(request CreateScheduleRequest) (CreateScheduleResult, error) {
	// Validation for time interval-based schedules
	if request.Attributes != nil && request.Attributes["frequency"] == "Time Interval" {
		if request.Minutes == nil || *request.Minutes < 5 {
			return CreateScheduleResult{},
				fmt.Errorf("time interval schedules require minutes >= 5")
		}
		if request.ScheduleInfo != nil {
			if retentionRaw, ok := request.ScheduleInfo["retention"]; ok {
				retentionStr, ok := retentionRaw.(string)
				if ok {
					// Parse retention string as H:M:S
					var h, m, s int
					n, err := fmt.Sscanf(retentionStr, "%d:%d:%d", &h, &m, &s)
					if err != nil || n != 3 {
						return CreateScheduleResult{}, fmt.Errorf("invalid retention format, expected H:M:S")
					}
					totalMinutes := h*60 + m + s/60
					if totalMinutes < 5 {
						return CreateScheduleResult{}, fmt.Errorf("time interval schedules require retention >= 5 minutes")
					}
				}
			}
		}
	}
	// Validation for single vs group snapshot schedule
	if request.ScheduleInfo != nil {
		_, hasVolumeID := request.ScheduleInfo["volumeID"]
		_, hasVolumes := request.ScheduleInfo["volumes"]
		if hasVolumeID && hasVolumes {
			return CreateScheduleResult{}, fmt.Errorf("scheduleInfo must have only one of volumeID or volumes, not both")
		}
		if !hasVolumeID && !hasVolumes {
			return CreateScheduleResult{}, fmt.Errorf("scheduleInfo must have either volumeID (single) or volumes (group)")
		}
		if hasVolumes {
			volumes, ok := request.ScheduleInfo["volumes"].([]interface{})
			if ok && len(volumes) < 2 {
				return CreateScheduleResult{}, fmt.Errorf("group snapshot schedules require at least 2 volumes")
			}
		}
	}
	params := structs.Map(request)
	response, err := c.CallAPIMethod("CreateSchedule", params)
	if err != nil {
		log.Print("CreateSchedule request failed")
		return CreateScheduleResult{}, err
	}

	var result CreateScheduleResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from CreateSchedule")
		return CreateScheduleResult{}, err
	}

	return result, nil
}

// ModifyScheduleRequest is used to modify a snapshot schedule, including marking it for deletion
// Corresponds to SolidFire Element API ModifySchedule

type ModifyScheduleRequest struct {
	ScheduleID      int                    `structs:"scheduleID"`
	ToBeDeleted     *bool                  `structs:"toBeDeleted,omitempty"`
	Attributes      map[string]interface{} `structs:"attributes,omitempty"`
	Hours           *int                   `structs:"hours,omitempty"`
	Minutes         *int                   `structs:"minutes,omitempty"`
	Paused          *bool                  `structs:"paused,omitempty"`
	Recurring       *bool                  `structs:"recurring,omitempty"`
	RunNextInterval *bool                  `structs:"runNextInterval,omitempty"`
	ScheduleName    *string                `structs:"scheduleName,omitempty"`
	ScheduleType    *string                `structs:"scheduleType,omitempty"`
	ScheduleInfo    map[string]interface{} `structs:"scheduleInfo,omitempty"`
	StartingDate    *string                `structs:"startingDate,omitempty"`
	Monthdays       []int                  `structs:"monthdays,omitempty"`
	Weekdays        []map[string]int       `structs:"weekdays,omitempty"`
}

type ModifyScheduleResult struct {
	Schedule Schedule `json:"schedule"`
}

func (c *Client) ModifySchedule(request ModifyScheduleRequest) (ModifyScheduleResult, error) {
	// Validation for time interval-based schedules
	if request.Attributes != nil && request.Attributes["frequency"] == "Time Interval" {
		if request.Minutes == nil || *request.Minutes < 5 {
			return ModifyScheduleResult{},
				fmt.Errorf("time interval schedules require minutes >= 5")
		}
		if request.ScheduleInfo != nil {
			if retentionRaw, ok := request.ScheduleInfo["retention"]; ok {
				retentionStr, ok := retentionRaw.(string)
				if ok {
					// Parse retention string as H:M:S
					var h, m, s int
					n, err := fmt.Sscanf(retentionStr, "%d:%d:%d", &h, &m, &s)
					if err != nil || n != 3 {
						return ModifyScheduleResult{}, fmt.Errorf("invalid retention format, expected H:M:S")
					}
					totalMinutes := h*60 + m + s/60
					if totalMinutes < 5 {
						return ModifyScheduleResult{}, fmt.Errorf("time interval schedules require retention >= 5 minutes")
					}
				}
			}
		}
	}
	// Validation for single vs group snapshot schedule
	if request.ScheduleInfo != nil {
		_, hasVolumeID := request.ScheduleInfo["volumeID"]
		_, hasVolumes := request.ScheduleInfo["volumes"]
		if hasVolumeID && hasVolumes {
			return ModifyScheduleResult{}, fmt.Errorf("scheduleInfo must have only one of volumeID or volumes, not both")
		}
		if !hasVolumeID && !hasVolumes {
			return ModifyScheduleResult{}, fmt.Errorf("scheduleInfo must have either volumeID (single) or volumes (group)")
		}
		if hasVolumes {
			volumes, ok := request.ScheduleInfo["volumes"].([]interface{})
			if ok && len(volumes) < 2 {
				return ModifyScheduleResult{}, fmt.Errorf("group snapshot schedules require at least 2 volumes")
			}
		}
	}
	params := structs.Map(request)
	response, err := c.CallAPIMethod("ModifySchedule", params)
	if err != nil {
		log.Print("ModifySchedule request failed")
		return ModifyScheduleResult{}, err
	}

	var result ModifyScheduleResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ModifySchedule")
		return ModifyScheduleResult{}, err
	}

	return result, nil
}
