package util

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
)

// SendCode 发送验证码邮件
func SendCode(to, code string) error {
	e := email.NewEmail()
	e.From = "gozknight <2690539295@qq.com>"
	e.To = []string{to}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("你的验证码是:<b>" + code + "</b>")
	//err := e.Send("smtp.qq.com:465", smtp.PlainAuth("", "2690539295@qq.com", "tihwlmpmvfkpddjb", "smtp.qq.com"))
	// 返回EOF时，关闭SSL重试
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "2690539295@qq.com", "tihwlmpmvfkpddjb", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true,
			ServerName: "smtp.qq.com",
		})
	if err != nil {
		return err
	}
	return nil
}
