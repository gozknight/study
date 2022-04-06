package test

import (
	"bytes"
	"cloud-disk/models"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
	"xorm.io/xorm"
)

func TestXorm(t *testing.T) {
	dburl := "root:root@tcp(127.0.0.1:3306)/cloud-disk?charset=utf8"
	engine, err := xorm.NewEngine("mysql", dburl)
	if err != nil {
		log.Fatalln(err)
	}
	user := make([]*models.UserBasic, 0)
	err = engine.Find(&user)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}
	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", "")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(dst.String())
}
