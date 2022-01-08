package schedule

import (
	"fmt"
	"github.com/ulstu-schedule/parser/schedule"
	"regexp"
	"strings"
)

var (
	KEIGroupPattern   = regexp.MustCompile(`^[А-Я]+[сдо]+-\d+$`)
	groupCharsPattern = regexp.MustCompile(`(?i)[а-я\d-]+`)
)

// GetDayGroupSchedule returns the text of the group schedule for the day taken from the UlSTU website.
func GetDayGroupSchedule(userGroupName string, userMsg string) (string, error) {
	if userMsg == "3" || userMsg == "сегодня" {
		return schedule.GetTextDayGroupSchedule(userGroupName, 0)
	} else {
		return schedule.GetTextDayGroupSchedule(userGroupName, 1)
	}
}

// ParseDayGroupSchedule returns the text of the group schedule for the day taken from the database table with backups
// of the group schedule.
func ParseDayGroupSchedule(scheduleJSON []byte, scheduleUpdateTimeFmt, userGroupName, userMsg string) (string, error) {
	fullSchedule, err := unmarshalFullSchedule(scheduleJSON)
	if err != nil {
		return "", err
	}

	if userMsg == "3" || userMsg == "сегодня" {
		todaySchedule, err := schedule.ParseDayGroupSchedule(fullSchedule, userGroupName, 0)
		if err != nil {
			return "", err
		}

		todayScheduleText := fmt.Sprintf("&#10071; По состоянию на %s\n\n", scheduleUpdateTimeFmt)
		todayScheduleText += schedule.ConvertDayGroupScheduleToText(todaySchedule, userGroupName, 0)

		return todayScheduleText, nil
	} else {
		tomorrowSchedule, err := schedule.ParseDayGroupSchedule(fullSchedule, userGroupName, 1)
		if err != nil {
			return "", err
		}

		tomorrowScheduleText := fmt.Sprintf("&#10071; По состоянию на %s\n\n", scheduleUpdateTimeFmt)
		tomorrowScheduleText += schedule.ConvertDayGroupScheduleToText(tomorrowSchedule, userGroupName, 1)

		return tomorrowScheduleText, nil
	}
}

// GetWeekGroupSchedule returns the path to the image with the week group schedule taken from the UlSTU website.
func GetWeekGroupSchedule(userGroupName, userMsg string) (caption string, weekSchedulePath string, err error) {
	if userMsg == "5" || userMsg == "текущая неделя" {
		weekSchedulePath, err = schedule.GetCurrWeekGroupScheduleImg(userGroupName)
		if err != nil {
			return
		}

		caption = fmt.Sprintf("Расписание %s на текущую неделю &#128071;\n\n", userGroupName)
	} else {
		weekSchedulePath, err = schedule.GetNextWeekGroupScheduleImg(userGroupName)
		if err != nil {
			return
		}

		caption = fmt.Sprintf("Расписание %s на следующую неделю &#128071;", userGroupName)
	}
	return
}

// ParseWeekGroupSchedule returns the path to the image with the week group schedule taken from the database table with
// backups of the group schedule.
func ParseWeekGroupSchedule(scheduleJSON []byte, scheduleUpdateTimeFmt, userGroupName, userMsg string) (string, string, error) {
	fullSchedule, err := unmarshalFullSchedule(scheduleJSON)
	if err != nil {
		return "", "", err
	}

	if userMsg == "5" || userMsg == "текущая неделя" {
		currWeekSchedule, err := schedule.ParseCurrWeekGroupSchedule(fullSchedule, userGroupName)
		if err != nil {
			return "", "", err
		}

		weekSchedulePath, err := schedule.ParseCurrWeekGroupScheduleImg(currWeekSchedule, userGroupName)
		if err != nil {
			return "", "", err
		}

		caption := fmt.Sprintf("&#10071; По состоянию на %s\n\nРасписание %s на текущую неделю &#128071;\n\n", scheduleUpdateTimeFmt, userGroupName)
		return caption, weekSchedulePath, nil
	} else {
		nextWeekSchedule, err := schedule.ParseNextWeekGroupSchedule(fullSchedule, userGroupName)
		if err != nil {
			return "", "", err
		}

		weekSchedulePath, err := schedule.ParseNextWeekGroupScheduleImg(nextWeekSchedule, userGroupName)
		if err != nil {
			return "", "", err
		}

		caption := fmt.Sprintf("&#10071; По состоянию на %s\n\nРасписание %s на следующую неделю &#128071;", scheduleUpdateTimeFmt, userGroupName)
		return caption, weekSchedulePath, nil
	}
}

func IsKEIGroup(groupName string) bool {
	return KEIGroupPattern.MatchString(groupName) && !strings.Contains(groupName, "РОНд")
}

// IsGroupReserver checks whether s is a group by accessing the database tables with backups of the group schedule.
func IsGroupReserver(groups []string, s string) (bool, string) {
	loweredS := strings.ToLower(s)
	convertedS := convertToGroupName(loweredS)

	for _, group := range groups {
		if convertedS == strings.ToLower(group) {
			return true, group
		}
	}

	return false, ""
}

// IsGroupParser checks whether s is a group by referring to the UlSTU website.
func IsGroupParser(s string) (bool, string) {
	loweredS := strings.ToLower(s)
	convertedS := convertToGroupName(loweredS)

	groups := schedule.GetGroups()
	for _, group := range groups {
		if strings.Contains(group, ", ") {
			splitGroups := strings.Split(group, ", ")
			for _, splitGroup := range splitGroups {
				if convertedS == strings.ToLower(splitGroup) {
					return true, splitGroup
				}
			}
		} else {
			if convertedS == strings.ToLower(group) {
				return true, group
			}
		}
	}

	return false, ""
}

func convertToGroupName(s string) string {
	sCleared := deleteExcessSymbols(s)
	groupNameInRunes := make([]rune, 0, len(sCleared)+2)

	var afterNum, quantityNum int
	for _, character := range sCleared {
		switch {
		case character == '-':
			afterNum, quantityNum = 0, 0
		case 48 <= character && character <= 57 && afterNum == 1 && quantityNum != 2:
			afterNum = 0
			quantityNum++
			groupNameInRunes = append(groupNameInRunes, '-')
		case 48 <= character && character <= 57 && afterNum == 0 && quantityNum != 2:
			quantityNum++
		case quantityNum == 2:
			quantityNum = 0
			groupNameInRunes = append(groupNameInRunes, '-')
		default:
			afterNum = 1
		}
		groupNameInRunes = append(groupNameInRunes, character)
	}

	return string(groupNameInRunes)
}

func deleteExcessSymbols(s string) string {
	results := groupCharsPattern.FindAllString(s, -1)
	return strings.Join(results, "")
}
