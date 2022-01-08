package postgres

import (
	"database/sql"
	"fmt"
	"github.com/ulstu-schedule/bot-vk/internal/model"
	"github.com/ulstu-schedule/bot-vk/internal/store"
)

const (
	studentsRepoName      = "vk_students"
	groupScheduleRepoName = "groups_schedule"
)

var (
	_ store.StudentRepository         = (*StudentRepository)(nil)
	_ store.GroupScheduleRepository   = (*GroupScheduleRepository)(nil)
	_ store.TeacherScheduleRepository = (*TeacherScheduleRepository)(nil)
)

type StudentRepository struct {
	store *StudentStore
}

func (r *StudentRepository) GetAllStudents() ([]model.Student, error) {
	students := []model.Student{}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", studentsRepoName)
	err := r.store.db.Select(&students, query)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *StudentRepository) Information(firstName, lastName string, userID int, groupName string, facultyID byte) error {
	student, err := r.GetStudent(userID)
	if err != nil {
		return err
	}

	if student != nil {
		r.UpdateStudent(firstName, lastName, userID, groupName, facultyID)
	} else {
		r.AddStudent(firstName, lastName, userID, groupName, facultyID)
	}
	return nil
}

func (r *StudentRepository) AddStudent(firstName, lastName string, userID int, groupName string, facultyID byte) {
	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, user_id, group_name, faculty_id) VALUES ($1, $2, $3, $4, $5)", studentsRepoName)
	r.store.db.MustExec(query, firstName, lastName, userID, groupName, facultyID)
}

func (r *StudentRepository) GetStudent(userID int) (*model.Student, error) {
	student := model.Student{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", studentsRepoName)
	err := r.store.db.Get(&student, query, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// there no student in the database table
	if student.ID == 0 {
		return nil, nil
	}

	return &student, nil
}

func (r *StudentRepository) UpdateStudent(firstName, lastName string, userID int, newGroupName string, newFacultyID byte) {
	query := fmt.Sprintf("UPDATE %s SET first_name=$2, last_name=$3, group_name=$4, faculty_id=$5 WHERE user_id=$1", studentsRepoName)
	r.store.db.MustExec(query, userID, firstName, lastName, newGroupName, newFacultyID)
}

type GroupScheduleRepository struct {
	store *ScheduleStore
}

func (r *GroupScheduleRepository) GetSchedule(groupName string) (*model.GroupSchedule, error) {
	schedule := model.GroupSchedule{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE group_name=$1", groupScheduleRepoName)
	err := r.store.db.Get(&schedule, query, groupName)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// if the group schedule is not in the database
	if schedule.Name == "" {
		return nil, nil
	}

	return &schedule, nil
}

func (r *GroupScheduleRepository) GetGroups() ([]string, error) {
	groups := []string{}
	query := fmt.Sprintf("SELECT group_name FROM %s", groupScheduleRepoName)
	err := r.store.db.Select(&groups, query)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

type TeacherScheduleRepository struct {
	// TODO: сделать по примеру GroupScheduleRepository
}
