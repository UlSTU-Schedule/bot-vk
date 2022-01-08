package vk

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/object"
)

func (b *Bot) getMainMenuKeyboard() *object.MessagesKeyboard {
	k := object.NewMessagesKeyboard(false)

	k.AddRow()
	k.AddTextButton("Моё расписание", "", "positive")
	k.AddTextButton("Изменить группу", "", "primary")

	return k
}

func (b *Bot) getScheduleMenuKeyboard() *object.MessagesKeyboard {
	k := object.NewMessagesKeyboard(false)

	k.AddRow()
	k.AddTextButton("Сегодня", "", "positive")
	k.AddTextButton("Завтра", "", "positive")

	k.AddRow()
	k.AddTextButton("Текущая неделя", "", "positive")
	k.AddTextButton("Следующая неделя", "", "positive")

	k.AddRow()
	k.AddTextButton("Назад", "", "primary")

	return k
}

func (b *Bot) hideKeyboard() *object.MessagesKeyboard {
	return object.NewMessagesKeyboard(false)
}

func (b *Bot) getDonationsInlineKeyboard() *object.MessagesKeyboard {
	k := object.NewMessagesKeyboardInline()

	k.AddRow()
	k.AddVKAppsButton(6471849, b.bot.GroupID, "", "Поддержать проект", "")

	return k
}

func (b *Bot) getSubscribingInlineKeyboard() *object.MessagesKeyboard {
	k := object.NewMessagesKeyboardInline()

	k.AddRow()
	link := fmt.Sprintf("https://vk.com/widget_community.php?act=a_subscribe_box&oid=-%d&state=1", b.bot.GroupID)
	k.AddOpenLinkButton(link, "Подписаться", "")

	return k
}

func (b *Bot) getCancelKeyboard() *object.MessagesKeyboard {
	k := object.NewMessagesKeyboard(false)

	k.AddRow()
	k.AddTextButton("Отмена", "", "primary")

	return k
}
