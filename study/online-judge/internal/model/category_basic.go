package model

import "gorm.io/gorm"

type CategoryBasic struct {
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name     string `gorm:"column:name;type:varchar(100);" json:"name"`
	ParentId string `gorm:"column:parent_id;type:int(11);" json:"parent_id"`
	gorm.Model
}

func (c *CategoryBasic) TableName() string {
	return "category_basic"
}
