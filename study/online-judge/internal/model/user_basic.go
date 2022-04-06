package model

import "gorm.io/gorm"

type UserBasic struct {
	Identity         string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name             string `gorm:"column:name;type:varchar(100);" json:"name"`
	Password         string `gorm:"column:password;type:varchar(32);" json:"password"`
	Phone            string `gorm:"column:phone;type:varchar(20);" json:"phone"`
	Email            string `gorm:"column:email;type:varchar(100);" json:"email"`
	FinishProblemNum int64  `gorm:"column:finish_problem_num;type:int(11);" json:"finish_problem_num"`
	SubmitNum        int64  `gorm:"column:submit_num;type:int(11);" json:"submit_num"`
	gorm.Model
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
