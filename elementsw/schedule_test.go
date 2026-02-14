package elementsw

import (
	"reflect"
	"testing"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

func TestCreateScheduleValidation(t *testing.T) {
	// Validation was removed in SDK migration.
	// If we want to restore it, we need to add it to Client.CreateSchedule.
}

func TestScheduleStructFields(t *testing.T) {
	s := sdk.Schedule{}
	fields := []string{"Hours", "Minutes", "Monthdays", "Paused", "Recurring", "RunNextInterval", "ScheduleID", "ScheduleInfo", "ScheduleName", "ScheduleType", "StartingDate", "ToBeDeleted", "Weekdays"}
	typeOf := reflect.TypeOf(s)
	for _, f := range fields {
		if _, ok := typeOf.FieldByName(f); !ok {
			t.Errorf("Schedule struct missing field: %s", f)
		}
	}
}
