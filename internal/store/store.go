package store

// ScheduleStore represents a database with backups of schedules.
type ScheduleStore interface {
	// GroupSchedule represents a database table with backups of group schedules.
	GroupSchedule() *GroupScheduleRepository

	// TeacherSchedule represents a database table with backups of teacher schedules.
	TeacherSchedule() *TeacherScheduleRepository
}

// StudentStore represents a database with bot users in Telegram.
type StudentStore interface {
	// Student represents a table in the database with users.
	Student() *StudentRepository
}
