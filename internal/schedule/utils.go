package schedule

import (
	"github.com/mailru/easyjson"
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
