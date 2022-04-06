package util

import (
	"math/rand"
	"strconv"
	"time"
)

// GetRandomCode 随机生成验证码
func GetRandomCode() (ans string) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		t := rand.Intn(10)
		ans += strconv.Itoa(t)
	}
	return
}
