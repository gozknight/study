package test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gozknight.com/online-judge/internal/model"
	"testing"
)

func TestGet(t *testing.T) {
	url := "root:root@tcp(127.0.0.1:3306)/online-judge?charset=utf8&parseTime=true&loc=Local"
	orm, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	data := make([]*model.UserBasic, 0)
	err = orm.Find(&data).Error
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range data {
		fmt.Printf("User ==> %v\n", u)
	}
}
