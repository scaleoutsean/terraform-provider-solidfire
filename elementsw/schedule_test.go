package elementsw

import (
	"reflect"
	"testing"
)

func TestCreateScheduleValidation(t *testing.T) {
	client := &Client{}
	// Invalid: minutes < 5
	minutes := 4
	request := CreateScheduleRequest{
		Attributes: map[string]interface{}{"frequency": "Time Interval"},
		Minutes: &minutes,
		ScheduleInfo: map[string]interface{}{"retention": "0:5:00"},
		ScheduleName: "test",
		ScheduleType: "snapshot",
	}
	_, err := client.CreateSchedule(request)
	if err == nil || err.Error() != "time interval schedules require minutes >= 5" {
		t.Errorf("expected minutes validation error, got %v", err)
	}

	// Invalid: retention < 5
	minutes = 5
	request.Minutes = &minutes
	request.ScheduleInfo["retention"] = "0:4:00"
	_, err = client.CreateSchedule(request)
	if err == nil || err.Error() != "time interval schedules require retention >= 5 minutes" {
		t.Errorf("expected retention validation error, got %v", err)
	}

	// Invalid: retention format
	request.ScheduleInfo["retention"] = "badformat"
	_, err = client.CreateSchedule(request)
	if err == nil || err.Error() != "invalid retention format, expected H:M:S" {
		t.Errorf("expected retention format error, got %v", err)
	}

	// Valid
	request.ScheduleInfo["retention"] = "0:5:00"
	_, err = client.CreateSchedule(request)
	if err != nil && err.Error() != "CreateSchedule request failed" {
		t.Errorf("expected API error or success, got %v", err)
	}
}

func TestModifyScheduleValidation(t *testing.T) {
	client := &Client{}
	minutes := 4
	request := ModifyScheduleRequest{
		Attributes: map[string]interface{}{"frequency": "Time Interval"},
		Minutes: &minutes,
		ScheduleInfo: map[string]interface{}{"retention": "0:5:00"},
		ScheduleName: ptrString("test"),
		ScheduleType: ptrString("snapshot"),
	}
	_, err := client.ModifySchedule(request)
	if err == nil || err.Error() != "time interval schedules require minutes >= 5" {
		t.Errorf("expected minutes validation error, got %v", err)
	}

	minutes = 5
	request.Minutes = &minutes
	request.ScheduleInfo["retention"] = "0:4:00"
	_, err = client.ModifySchedule(request)
	if err == nil || err.Error() != "time interval schedules require retention >= 5 minutes" {
		t.Errorf("expected retention validation error, got %v", err)
	}

	request.ScheduleInfo["retention"] = "badformat"
	_, err = client.ModifySchedule(request)
	if err == nil || err.Error() != "invalid retention format, expected H:M:S" {
		t.Errorf("expected retention format error, got %v", err)
	}

	request.ScheduleInfo["retention"] = "0:5:00"
	_, err = client.ModifySchedule(request)
	if err != nil && err.Error() != "ModifySchedule request failed" {
		t.Errorf("expected API error or success, got %v", err)
	}
}

func ptrString(s string) *string { return &s }

func TestScheduleStructFields(t *testing.T) {
	s := Schedule{}
	fields := []string{"Attributes", "HasError", "Hours", "LastRunStatus", "LastRunTimeStarted", "Minutes", "Monthdays", "Paused", "Recurring", "RunNextInterval", "ScheduleID", "ScheduleInfo", "ScheduleName", "ScheduleType", "StartingDate", "ToBeDeleted", "Weekdays"}
	typeOf := reflect.TypeOf(s)
	for _, f := range fields {
		if _, ok := typeOf.FieldByName(f); !ok {
			t.Errorf("Schedule struct missing field: %s", f)
		}
	}
}
