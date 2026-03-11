package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	FirstName            string       `gorm:"column:first_name;size:100;not null"`
	LastName             string       `gorm:"column:last_name;size:100;not null"`
	Email                string       `gorm:"column:email;size:255;uniqueIndex"`
}

type Subject struct {
	gorm.Model
	Name                 string       `gorm:"column:name;size:200;not null"`
}

type Lesson struct {
	gorm.Model
	Date                 time.Time    `gorm:"column:date;not null"`
	SubjectID            int          `gorm:"column:subject_id"`
	Subject              Subject      `gorm:"foreignKey:SubjectID"`
}

type Attendance struct {
	gorm.Model
	LessonID             int          `gorm:"column:lesson_id"`
	Lesson               Lesson       `gorm:"foreignKey:LessonID"`
	StudentID            int          `gorm:"column:student_id"`
	Student              Student      `gorm:"foreignKey:StudentID"`
	Present              bool         `gorm:"column:present;default:FALSE"`
}
