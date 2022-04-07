package model

import "gorm.io/gorm"

type TestCase struct {
	Identity        string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	Input           string `gorm:"column:input;type:text;" json:"input"`
	Output          string `gorm:"column:output;type:text;" json:"output"`
	gorm.Model
}

func (t *TestCase) TableName() string {
	return "test_case"
}
