package postgres

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB(databaseUrl string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

type StudentStore struct {
	db                *sqlx.DB
	studentRepository *StudentRepository
}

func NewStudentStore(db *sqlx.DB) *StudentStore {
	return &StudentStore{
		db: db,
	}
}

func (s *StudentStore) Student() *StudentRepository {
	if s.studentRepository != nil {
		return s.studentRepository
	}

	s.studentRepository = &StudentRepository{
		store: s,
	}

	return s.studentRepository
}

type ScheduleStore struct {
	db                        *sqlx.DB
	groupScheduleRepository   *GroupScheduleRepository
	teacherScheduleRepository *TeacherScheduleRepository
}

func NewScheduleStore(db *sqlx.DB) *ScheduleStore {
	return &ScheduleStore{
		db: db,
	}
}

func (s *ScheduleStore) GroupSchedule() *GroupScheduleRepository {
	if s.groupScheduleRepository != nil {
		return s.groupScheduleRepository
	}

	s.groupScheduleRepository = &GroupScheduleRepository{
		store: s,
	}

	return s.groupScheduleRepository
}

func (s *ScheduleStore) TeacherSchedule() *TeacherScheduleRepository {
	// TODO: сделать по примеру GroupSchedule()
	panic("implement me!")
}
