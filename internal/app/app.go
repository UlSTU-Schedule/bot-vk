package app

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/ulstu-schedule/bot-vk/internal/config"
	"github.com/ulstu-schedule/bot-vk/internal/store/postgres"
	"github.com/ulstu-schedule/bot-vk/internal/vk"
	"log"
)

// Run runs the bot.
func Run(configsPath string) {
	cfg, err := config.New(configsPath)
	if err != nil {
		log.Fatal(err)
	}

	studentDB, err := postgres.NewDB(cfg.StudentDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	studentStore := postgres.NewStudentStore(studentDB)

	scheduleDB, err := postgres.NewDB(cfg.ScheduleDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	scheduleStore := postgres.NewScheduleStore(scheduleDB)

	newVK := api.NewVK(cfg.Token)

	group, err := newVK.GroupsGetByID(api.Params{})
	if err != nil {
		log.Fatal(err)
	}

	lp, err := longpoll.NewLongPoll(newVK, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := vk.NewBot(lp, cfg.StickerIDs, cfg.Messages, cfg.Commands, studentStore, scheduleStore, cfg.Faculties)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.RunPolling()
	if err != nil {
		log.Fatal(err)
	}
}
