package store

import (
	"github.com/ulstu-schedule/bot-vk/internal/model"
)

// StudentRepository represents a table in a database with bot users in Telegram.
type StudentRepository interface {
	// GetAllStudents returns a slice of all bot users.
	GetAllStudents() ([]model.Student, error)

	// Information executes AddStudent if there is no user in the table and executes UpdateStudent if there is in the table.
	Information(firstName, lastName string, userID int, groupName string, facultyID byte) error

	// AddStudent adds the user to the database table with the bot users.
	AddStudent(firstName, lastName string, userID int, groupName string, facultyID byte)

	// GetStudent returns information about the bot user.
	GetStudent(userID int) (*model.Student, error)

	// UpdateStudent updates information about the bot user.
	UpdateStudent(firstName, lastName string, userID int, newGroupName string, facultyID byte)
}

// GroupScheduleRepository represents a database table with backups of group schedules.
type GroupScheduleRepository interface {
	// GetSchedule returns a backup copy of the group schedule with additional information.
	GetSchedule(groupName string) (*model.GroupSchedule, error)

	// GetGroups returns a slice of all groups that are in the database with backups of groups schedules.
	GetGroups() ([]string, error)
}

// TeacherScheduleRepository represents a database table with backups of teacher schedules.
type TeacherScheduleRepository interface {
	// TODO: сделать по примеру с группами
}
