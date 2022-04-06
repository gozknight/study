package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "gozknight <2690539295@qq.com>"
	e.To = []string{"gozknight@163.com"}
	e.Subject = "验证码"
	e.HTML = []byte("你的验证码是:<h1>123</h1>")
	//err := e.Send("smtp.qq.com:465", smtp.PlainAuth("", "2690539295@qq.com", "tihwlmpmvfkpddjb", "smtp.qq.com"))
	// 返回EOF时，关闭SSL重试
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "2690539295@qq.com", "tihwlmpmvfkpddjb", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true,
			ServerName: "smtp.qq.com",
		})
	if err != nil {
		t.Fatal(err)
	}
}
