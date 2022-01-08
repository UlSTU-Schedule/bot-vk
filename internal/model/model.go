package model

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

// Student represents the bot user in VK.
type Student struct {
	ID        int
	UserID    int    `db:"user_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	GroupName string `db:"group_name"`
	FacultyID byte   `db:"faculty_id"`
}

// GroupSchedule represents an entry in the database table with backups of the group schedule.
type GroupSchedule struct {
	ID         int
	Name       string         `db:"group_name"`
	UpdateTime time.Time      `db:"update_time"`
	Info       types.JSONText `db:"info"`
}

// TeacherSchedule represents an entry in the database table with backups of the teacher schedule.
type TeacherSchedule struct {
	// TODO: сделать по аналогии с GroupSchedule
}
