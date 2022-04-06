package util

import (
	"crypto/md5"
	"fmt"
)

// MD5 加密密码
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
