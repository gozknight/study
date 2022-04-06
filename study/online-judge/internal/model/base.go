package model

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var ORM = Init()

func Init() *gorm.DB {
	url := "root:root@tcp(127.0.0.1:3306)/online-judge?charset=utf8&parseTime=true&loc=Local"
	orm, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	return orm
}
