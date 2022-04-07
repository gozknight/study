package model

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	Identity          string             `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`
	Title             string             `gorm:"column:title;type:varchar(255)" json:"title"`
	Content           string             `gorm:"column:content;type:text" json:"content"`
	MaxMemory         string             `gorm:"column:max_memory;type:int(11)" json:"max_memory"`
	MaxRuntime        string             `gorm:"column:max_runtime;type:int(11)" json:"max_runtime"`
	TestCase          []*TestCase        `gorm:"foreignKey:problem_identity;references:identity" json:"test_case"`
	gorm.Model
}

func (p *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	key := "%" + keyword + "%"
	tx := ORM.Model(new(ProblemBasic)).
		Preload("ProblemCategories").
		Preload("ProblemCategories.CategoryBasic").
		Where("title like ? or content like ?", key, key)
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ?)", categoryIdentity)
	}
	return tx
}
