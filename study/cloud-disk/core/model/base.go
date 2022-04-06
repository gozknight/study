package model

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"xorm.io/xorm"
)

var Engine = Init()
var once *sync.Once

func Init() *xorm.Engine {
	dburl := "root:root@tcp(127.0.0.1:3306)/cloud-disk?charset=utf8"
	engine, err := xorm.NewEngine("mysql", dburl)
	if err != nil {
		log.Fatalln(err)
	}
	return engine
}
