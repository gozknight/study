package test

import (
	"encoding/base64"
	"fmt"
	"gozknight.com/online-judge/internal/util"
	"testing"
)

func TestMD5(t *testing.T) {
	c := base64.StdEncoding.EncodeToString([]byte("i miss you"))
	fmt.Println(c)
	code := util.MD5("i miss you")
	fmt.Println(code)
}
