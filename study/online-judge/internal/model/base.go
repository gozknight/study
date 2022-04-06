package model

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var ORM = InitMySQL()
var RDB = InitRedis()

func InitRedis() *redis.Client {
	var rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
func InitMySQL() *gorm.DB {
	url := "root:root@tcp(127.0.0.1:3306)/online-judge?charset=utf8&parseTime=true&loc=Local"
	orm, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	return orm
}
