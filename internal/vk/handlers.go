package vk

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/ulstu-schedule/bot-vk/internal/schedule"
	"github.com/ulstu-schedule/parser/types"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

func (b *Bot) handleNewMessages() {
	b.bot.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		go func(obj events.MessageNewObject) {
			switch {
			case obj.Message.Text != "":
				if err := b.handleTextMsg(obj.Message.Text, obj.Message.FromID); err != nil {
					b.handleError(err, obj)
				}
			case len(obj.Message.Attachments) > 0:
				if err := b.handleAttachment(obj.Message.Attachments[0].Type, obj.Message.FromID); err != nil {
					b.handleError(err, obj)
				}
			}
		}(obj)
	})
}

func (b *Bot) handleError(err error, obj events.MessageNewObject) {
	userID := obj.Message.FromID
	userMsg := obj.Message.Text

	log.Printf("[VK] id%d: text=%s, attachments=%v", userID, userMsg, obj.Message.Attachments)
	log.Printf("[VK] ERROR: %s", err)

	_ = b.markAsImportant(userID)

	switch err.(type) {
	case *types.UnavailableScheduleError, *types.IncorrectLinkError:
		_ = b.sendMessage(userID, b.messages.ScheduleIsUnavailable, "", nil)
	case *types.StatusCodeError:
		_ = b.sendMessage(userID, b.messages.ServerError, "", nil)
	case *types.IncorrectDateError:
		_ = b.sendMessage(userID, b.messages.IncorrectDateError, "", nil)
	default:
		_ = b.sendMessage(userID, b.messages.UnknownError, "", nil)
	}
}

func (b *Bot) handleAttachment(attachmentType string, userID int) error {
	switch attachmentType {
	case "sticker":
		return b.handleSticker(userID)
	case "audio_message":
		return b.handleAudioMsg(userID)
	default:
		return nil
	}
}

func (b *Bot) handleSticker(userID int) error {
	stickerIDs := b.stickerIDs.ToSticker
	rndStickerID := stickerIDs[rand.Intn(len(stickerIDs))]
	return b.sendSticker(userID, rndStickerID)
}

func (b *Bot) handleAudioMsg(userID int) error {
	stickerIDs := b.stickerIDs.ToAudioMessage
	rndStickerID := stickerIDs[rand.Intn(len(stickerIDs))]
	return b.sendSticker(userID, rndStickerID)
}

func (b *Bot) handleTextMsg(userMsg string, userID int) error {
	userMsgLowered := strings.ToLower(userMsg)

	switch {
	case contains(b.commands.Whole.GetScheduleForDay, userMsgLowered):
		return b.handleGetScheduleForDayMsg(userID, userMsgLowered)
	case contains(b.commands.Whole.GetScheduleForWeek, userMsgLowered) ||
		containsPartial(b.commands.Partial.GetScheduleForWeek, userMsgLowered):
		return b.handleGetScheduleForWeekMsg(userID, userMsgLowered)
	case contains(b.commands.Whole.ChangeGroup, userMsgLowered) ||
		containsPartial(b.commands.Partial.ChangeGroup, userMsgLowered):
		return b.handleChangeGroupMsg(userID)
	case contains(b.commands.Whole.BackToStartMenu, userMsgLowered) ||
		containsPartial(b.commands.Partial.BackToStartMenu, userMsgLowered):
		return b.handleBackToStartMenuMsg(userID)
	case contains(b.commands.Whole.GoToScheduleMenu, userMsgLowered) ||
		containsPartial(b.commands.Partial.GoToScheduleMenu, userMsgLowered):
		return b.handleGoToScheduleMenuMsg(userID)
	case contains(b.commands.Whole.Start, userMsgLowered) ||
		containsPartial(b.commands.Partial.Start, userMsgLowered):
		return b.handleStartMsg(userID)
	case contains(b.commands.Whole.Session, userMsgLowered) ||
		containsPartial(b.commands.Partial.Session, userMsgLowered):
		return b.handleSessionMsg(userID)
	case containsPartial(b.commands.Partial.ExpressGratitude, userMsgLowered):
		return b.handleExpressGratitudeMsg(userID)
	default:
		return b.handleUnknownMsg(userID, userMsg)
	}
}

func (b *Bot) handleStartMsg(userID int) error {
	student, err := b.studentStore.Student().GetStudent(userID)
	if err != nil {
		return err
	}

	if student != nil {
		return b.sendMessage(userID, b.messages.StartWithGroup, "", b.getMainMenuKeyboard())
	} else {
		return b.sendMessage(userID, b.messages.StartWithoutGroup, "", b.hideKeyboard())
	}
}

func (b *Bot) handleChangeGroupMsg(userID int) error {
	return b.sendMessage(userID, b.messages.ChangeGroup, "", b.getCancelKeyboard())
}

func (b *Bot) handleGetScheduleForDayMsg(userID int, userMsg string) error {
	student, err := b.studentStore.Student().GetStudent(userID)
	if err != nil {
		return err
	}

	if student != nil {
		daySchedule, err := schedule.GetDayGroupSchedule(student.GroupName, userMsg)
		if err != nil {
			groupScheduleJSON, err := b.scheduleStore.GroupSchedule().GetSchedule(student.GroupName)
			if err != nil {
				return err
			}

			updateTimeFmt := groupScheduleJSON.UpdateTime.Format("15:04:05 02.01.2006")

			daySchedule, err = schedule.ParseDayGroupSchedule(groupScheduleJSON.Info, updateTimeFmt, student.GroupName, userMsg)
			if err != nil {
				return err
			}
		}

		if schedule.IsKEIGroup(student.GroupName) {
			daySchedule += b.messages.ChangesInKEISchedule
		}

		return b.sendMessage(userID, daySchedule, "", b.getScheduleMenuKeyboard())
	} else {
		return b.sendMessage(userID, b.messages.GroupNotSelected, "", b.hideKeyboard())
	}
}

func (b *Bot) handleGetScheduleForWeekMsg(userID int, userMsg string) error {
	isSubscriber, err := b.isSubscriber(userID)
	if err != nil {
		return err
	}

	if isSubscriber {
		student, err := b.studentStore.Student().GetStudent(userID)
		if err != nil {
			return err
		}

		if student != nil {
			_ = b.sendMessage(userID, "Генерирую расписание &#9203;", "", b.getScheduleMenuKeyboard())

			answer, attachment, err := b.getWeeklySchedule(userID, student.GroupName, userMsg)
			if err != nil {
				return err
			}

			return b.sendMessage(userID, answer, attachment, b.getScheduleMenuKeyboard())
		} else {
			return b.sendMessage(userID, b.messages.GroupNotSelected, "", b.hideKeyboard())
		}
	} else {
		return b.sendMessage(userID, b.messages.StudentNotSubscribed, "", b.getSubscribingInlineKeyboard())
	}
}

// getWeeklySchedule takes lowered user message.
func (b *Bot) getWeeklySchedule(userID int, userGroup, userMsg string) (string, string, error) {
	caption, weekSchedulePath, err := schedule.GetWeekGroupSchedule(userGroup, userMsg)
	if err != nil {
		groupScheduleJSON, err := b.scheduleStore.GroupSchedule().GetSchedule(userGroup)
		if err != nil {
			return "", "", err
		}

		updateTimeFmt := groupScheduleJSON.UpdateTime.Format("15:04:05 02.01.2006")

		caption, weekSchedulePath, err = schedule.ParseWeekGroupSchedule(groupScheduleJSON.Info, updateTimeFmt, userGroup, userMsg)
		if err != nil {
			return "", "", err
		}
	}
	if weekSchedulePath != "" {
		defer os.Remove(weekSchedulePath)
	}

	if (userMsg == "5" || userMsg == "текущая неделя") && schedule.IsKEIGroup(userGroup) {
		caption += b.messages.ChangesInKEISchedule
	}

	attachment, err := b.uploadMessagesPhoto(userID, weekSchedulePath)
	if err != nil {
		return "", "", err
	}

	return caption, attachment, nil
}

func (b *Bot) handleBackToStartMenuMsg(userID int) error {
	student, err := b.studentStore.Student().GetStudent(userID)
	if err != nil {
		return err
	}

	if student != nil {
		return b.sendMessage(userID, b.messages.Back, "", b.getMainMenuKeyboard())
	} else {
		return b.sendMessage(userID, b.messages.ChangeGroup, "", b.hideKeyboard())
	}
}

func (b *Bot) handleGoToScheduleMenuMsg(userID int) error {
	student, err := b.studentStore.Student().GetStudent(userID)
	if err != nil {
		return err
	}

	if student != nil {
		answer := fmt.Sprintf("Твоя группа: %s &#128204; \n\n", student.GroupName)
		answer += b.messages.InfoWithGroup

		return b.sendMessage(userID, answer, "", b.getScheduleMenuKeyboard())
	} else {
		return b.sendMessage(userID, b.messages.InfoWithoutGroup, "", b.hideKeyboard())
	}
}

func (b *Bot) handleExpressGratitudeMsg(userID int) error {
	isSubscriber, err := b.isSubscriber(userID)
	if err != nil {
		return err
	}

	if isSubscriber {
		return b.sendMessage(userID, b.messages.Thanks, "", b.getDonationsInlineKeyboard())
	} else {
		return b.sendMessage(userID, b.messages.ThanksNotSubscribed, "", b.getSubscribingInlineKeyboard())
	}
}

func (b *Bot) handleSessionMsg(userID int) error {
	return b.sendMessage(userID, b.messages.Session, "", nil)
}

func (b *Bot) handleUnknownMsg(userID int, userMsg string) error {
	if isGroup, groupName := schedule.IsGroupParser(userMsg); isGroup {
		return b.updateGroup(userID, groupName)
	} else {
		groups, err := b.scheduleStore.GroupSchedule().GetGroups()
		if err != nil {
			return err
		}

		if isGroup, groupName = schedule.IsGroupReserver(groups, userMsg); isGroup {
			return b.updateGroup(userID, groupName)
		} else {
			return b.sendMessage(userID, b.messages.IncorrectInput, "", nil)
		}
	}
}

func (b *Bot) updateGroup(userID int, groupName string) error {
	firstName, lastName, err := b.getUserNames(userID)
	if err != nil {
		return err
	}

	facultyID := b.determineFacultyID(groupName)

	err = b.studentStore.Student().Information(firstName, lastName, userID, groupName, facultyID)
	if err != nil {
		return err
	}

	answer := fmt.Sprintf("Твоя группа обновлена на %s &#9989;\n\n", groupName)
	answer += b.messages.InfoWithGroup

	return b.sendMessage(userID, answer, "", b.getScheduleMenuKeyboard())
}

func (b *Bot) determineFacultyID(groupName string) byte {
	for _, faculty := range b.faculties {
		for _, group := range faculty.Groups {
			expr := fmt.Sprintf(`(?i)^%s[\d]+$`, group)
			groupRegexp := regexp.MustCompile(expr)
			if groupRegexp.MatchString(groupName) {
				return faculty.ID
			}
		}
	}
	if schedule.KEIGroupPattern.MatchString(groupName) {
		return 2
	}
	return 12
}

func containsPartial(s []string, e string) bool {
	amongRegexp := regexp.MustCompile(strings.Join(s, "|"))
	return amongRegexp.MatchString(e)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if e == a {
			return true
		}
	}
	return false
}
