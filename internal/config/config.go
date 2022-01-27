package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
)

const (
	answersConfigPath  = "answers"
	commandsConfigPath = "commands"
	facultiesJSONPath  = "configs/faculties.json"
)

type Config struct {
	Token               string
	StudentDatabaseURL  string
	ScheduleDatabaseURL string

	StickerIDs StickerIDs
	Messages   Messages
	Commands   Commands
	Faculties  []Faculty
}

// Commands represents the bot commands.
type Commands struct {
	Whole
	Partial
}

// Whole represents commands, which should be a whole message (е.g. msg='завтра', cmd='завтра').
type Whole struct {
	Start              []string `mapstructure:"start"`
	GoToScheduleMenu   []string `mapstructure:"go_to_schedule_menu"`
	ChangeGroup        []string `mapstructure:"change_group"`
	GetScheduleForDay  []string `mapstructure:"get_schedule_for_day"`
	GetScheduleForWeek []string `mapstructure:"get_schedule_for_week"`
	BackToStartMenu    []string `mapstructure:"back_to_start_menu"`
}

// Partial represents commands that can be inside a message (msg='можно моё расписание?', cmd='расписание').
type Partial struct {
	Start              []string `mapstructure:"start"`
	GoToScheduleMenu   []string `mapstructure:"go_to_schedule_menu"`
	GetScheduleForWeek []string `mapstructure:"get_schedule_for_week"`
	ChangeGroup        []string `mapstructure:"change_group"`
	BackToStartMenu    []string `mapstructure:"back_to_start_menu"`
	ExpressGratitude   []string `mapstructure:"express_gratitude"`
}

// Faculty represents UlSTU faculty.
type Faculty struct {
	Name   string
	ID     byte
	Groups []string
}

// StickerIDs represents the sticker IDs that the bot sends.
type StickerIDs struct {
	ToSticker      []int `mapstructure:"to_sticker"`
	ToAudioMessage []int `mapstructure:"to_audio_message"`
}

// Messages represents the messages that the bot sends to the user: regular, additional, and error.
type Messages struct {
	StartWithGroup       string `mapstructure:"start_with_group"`
	StartWithoutGroup    string `mapstructure:"start_without_group"`
	ChangeGroup          string `mapstructure:"change_group"`
	Back                 string `mapstructure:"back"`
	InfoWithoutGroup     string `mapstructure:"info_without_group"`
	InfoWithGroup        string `mapstructure:"info_with_group"`
	Thanks               string `mapstructure:"thanks"`
	IncorrectInput       string `mapstructure:"incorrect_input"`
	GroupNotSelected     string `mapstructure:"group_not_selected"`
	StudentNotSubscribed string `mapstructure:"student_not_subscribed"`

	ThanksNotSubscribed  string `mapstructure:"thanks_not_subscribed"`
	InfoNotSubscribed    string `mapstructure:"info_not_subscribed"`
	ChangesInKEISchedule string `mapstructure:"changes_in_kei_schedule"`

	ScheduleIsUnavailable string `mapstructure:"schedule_is_unavailable"`
	ServerError           string `mapstructure:"server_error"`
	UnknownError          string `mapstructure:"unknown_error"`
	IncorrectDateError    string `mapstructure:"incorrect_date_error"`
}

func New(configsPath string) (*Config, error) {
	cfg := &Config{}

	viper.AddConfigPath(configsPath)

	if err := fromAnswers(cfg); err != nil {
		return nil, err
	}

	if err := fromCommands(cfg); err != nil {
		return nil, err
	}

	if err := fromEnv(cfg); err != nil {
		return nil, err
	}

	if err := fromJson(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func fromAnswers(cfg *Config) error {
	if err := setUpViper(answersConfigPath); err != nil {
		return err
	}

	return unmarshalAnswers(cfg)
}

func fromCommands(cfg *Config) error {
	if err := setUpViper(commandsConfigPath); err != nil {
		return err
	}

	return unmarshalCommands(cfg)
}

func unmarshalCommands(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	return viper.UnmarshalKey("commands", &cfg.Commands)
}

func unmarshalAnswers(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("answers", &cfg.Messages); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("additions", &cfg.Messages); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("errors", &cfg.Messages); err != nil {
		return err
	}

	return viper.UnmarshalKey("sticker_ids", &cfg.StickerIDs)
}

func fromEnv(cfg *Config) error {
	if err := viper.BindEnv("TOKEN"); err != nil {
		return err
	}
	cfg.Token = viper.GetString("TOKEN")

	if err := viper.BindEnv("STUDENT_DATABASE_URL"); err != nil {
		return err
	}
	cfg.StudentDatabaseURL = viper.GetString("STUDENT_DATABASE_URL")

	if err := viper.BindEnv("SCHEDULE_DATABASE_URL"); err != nil {
		return err
	}
	cfg.ScheduleDatabaseURL = viper.GetString("SCHEDULE_DATABASE_URL")

	return nil
}

func fromJson(cfg *Config) error {
	data, err := ioutil.ReadFile(facultiesJSONPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cfg.Faculties)
}

func setUpViper(pathToConfigFile string) error {
	viper.SetConfigName(pathToConfigFile)
	return viper.ReadInConfig()
}
