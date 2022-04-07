package model

import "gorm.io/gorm"

type ProblemCategory struct {
	ProblemId     uint           `gorm:"column:problem_id;type:int(11)" json:"problem_id"`
	CategoryId    uint           `gorm:"column:category_id;type:int(11)" json:"category_id"`
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:category_id"`
	gorm.Model
}

func (pc *ProblemCategory) TableName() string {
	return "problem_category"
}
