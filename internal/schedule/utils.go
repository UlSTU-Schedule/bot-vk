package schedule

import (
	"github.com/mailru/easyjson"
	"github.com/ulstu-schedule/parser/schedule"
	"github.com/ulstu-schedule/parser/types"
)

func unmarshalFullSchedule(fullScheduleJSON []byte) (*types.Schedule, error) {
	fullSchedule := types.Schedule{}

	err := easyjson.Unmarshal(fullScheduleJSON, &fullSchedule)
	if err != nil {
		return nil, err
	}

	return &fullSchedule, nil
}

// GetSchoolWeekIdx returns index (0 or 1) of the school week depending on the LOWERED text command.
func GetSchoolWeekIdx(command string) int {
	additionalDays := 0

	switch command {
	case "4", "завтра":
		additionalDays = 1
		break
	case "6", "следующая неделя":
		additionalDays = 7
		break
	}

	weekNum, _ := schedule.GetWeekAndWeekDayNumbers(additionalDays)
	return weekNum
}
