package vk

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/ulstu-schedule/bot-vk/internal/config"
	"github.com/ulstu-schedule/bot-vk/internal/store/postgres"
	"log"
	"os"
	"strconv"
)

type Bot struct {
	bot           *longpoll.LongPoll
	stickerIDs    config.StickerIDs
	messages      config.Messages
	commands      config.Commands
	studentStore  *postgres.StudentStore
	scheduleStore *postgres.ScheduleStore
	faculties     []config.Faculty
}

func NewBot(lp *longpoll.LongPoll, stickerIDs config.StickerIDs, messages config.Messages, commands config.Commands, studentStore *postgres.StudentStore, scheduleStore *postgres.ScheduleStore, faculties []config.Faculty) (*Bot, error) {
	return &Bot{bot: lp, stickerIDs: stickerIDs, messages: messages, commands: commands, studentStore: studentStore, scheduleStore: scheduleStore, faculties: faculties}, nil
}

func (b *Bot) RunPolling() error {
	log.Println("The VK bot was launched!")

	b.handleNewMessages()
	return b.bot.Run()
}

func (b *Bot) getUserNames(userID int) (string, string, error) {
	response, err := b.bot.VK.UsersGet(api.Params{
		"user_ids": userID,
	})
	if err != nil {
		return "", "", err
	}

	return response[0].FirstName, response[0].LastName, nil
}

func (b *Bot) sendMessage(peerID int, text, attachment string, keyboard *object.MessagesKeyboard) error {
	if keyboard != nil {
		_, err := b.bot.VK.MessagesSend(api.Params{
			"user_id":    peerID,
			"random_id":  0,
			"attachment": attachment,
			"message":    text,
			"keyboard":   keyboard,
		})
		return err
	} else {
		_, err := b.bot.VK.MessagesSend(api.Params{
			"user_id":    peerID,
			"random_id":  0,
			"attachment": attachment,
			"message":    text,
		})
		return err
	}
}

func (b *Bot) sendSticker(peerID, stickerID int) error {
	_, err := b.bot.VK.MessagesSend(api.Params{
		"user_id":    peerID,
		"random_id":  0,
		"sticker_id": stickerID,
	})
	return err
}

func (b *Bot) isSubscriber(userID int) (bool, error) {
	isMemberInt, err := b.bot.VK.GroupsIsMember(api.Params{
		"group_id": strconv.Itoa(b.bot.GroupID),
		"user_id":  userID,
	})
	if isMemberInt == 1 {
		return true, err
	} else {
		return false, err
	}
}

func (b *Bot) uploadMessagesPhoto(peerID int, pathToPhoto string) (string, error) {
	response, err := os.Open(pathToPhoto)
	if err != nil {
		return "", err
	}
	defer response.Close()

	photo, err := b.bot.VK.UploadMessagesPhoto(peerID, response)
	if err != nil {
		return "", err
	}

	photoAttachment := fmt.Sprintf("photo%d_%d", photo[0].OwnerID, photo[0].ID)
	return photoAttachment, nil
}

func (b *Bot) markAsImportant(peerID int) error {
	_, err := b.bot.VK.MessagesMarkAsImportantConversation(api.Params{
		"peer_id":   peerID,
		"important": 1,
		"group_id":  b.bot.GroupID,
	})
	return err
}
