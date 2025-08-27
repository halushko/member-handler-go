package database

import (
	"log"

	"github.com/halushko/core-go/sqlite"
)

var DB sqlite.DBI

const (
	TStudents          = "students"
	TCourses           = "courses"
	TStudentsInCourses = "students_in_courses"

	CUserId     = "user_id"
	CLogin      = "login"
	CFirstName  = "first_name"
	CLastName   = "last_name"
	CMiddleName = "middle_name"
	CGroupId    = "group_id"

	CCourseId       = "course_id"
	CCourseFullName = "course_full_name"
	CChatId         = "chat_id"
	CSheetId        = "sheet_id"
)

var tblStudents = sqlite.Table{
	Name: TStudents,
	Columns: []sqlite.Column{
		{
			Name:     CUserId,
			Type:     sqlite.Integer,
			IsUnique: true,
		},
		{
			Name:     CLogin,
			Type:     sqlite.Text,
			IsUnique: true,
		},
		{
			Name: CFirstName,
			Type: sqlite.Text,
		},
		{
			Name: CLastName,
			Type: sqlite.Text,
		},
		{
			Name: CMiddleName,
			Type: sqlite.Text,
		},
		{
			Name: CGroupId,
			Type: sqlite.Text,
		},
	},
}
var tblCourses = sqlite.Table{
	Name: TCourses,
	Columns: []sqlite.Column{
		{
			Name:         CCourseId,
			Type:         sqlite.Integer,
			IsPrimaryKey: true,
		},
		{
			Name:      CChatId,
			Type:      sqlite.Text,
			IsUnique:  true,
			IsNotNull: true,
		},
		{
			Name: CCourseFullName,
			Type: sqlite.Text,
		},
		{
			Name:     CSheetId,
			Type:     sqlite.Text,
			IsUnique: true,
		},
	},
}
var tblStudentInCourse = sqlite.Table{
	Name: TStudentsInCourses,
	Columns: []sqlite.Column{
		{
			Name:      CUserId,
			Type:      sqlite.Integer,
			IsNotNull: true,
		},
		{
			Name:      CCourseId,
			Type:      sqlite.Integer,
			IsNotNull: true,
		},
	},
}

func Init() error {
	tables := []sqlite.Table{
		tblStudents,
		tblCourses,
		tblStudentInCourse,
	}

	project := sqlite.DBInfo{
		Name:    "course_chat",
		Project: "ist_chat_bot",
		Tables:  tables,
	}
	db, err := sqlite.Init(project)

	if err != nil {
		log.Printf("[ERROR] can't create course_chat: %v", err)
		return err
	}

	DB = db
	return nil
}
