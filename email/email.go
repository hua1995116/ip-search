package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

func sendMail() {
	auth := smtp.PlainAuth("", UserEmail, Mail_Password, Mail_Smtp_Host)
	to := []string{"qiufenghyf@163.com"}
	nickname := "秋风"
	user := UserEmail

	subject := "标题"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := "邮件内容."
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(Mail_Smtp_Host+Mail_Smtp_Port, auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
	}
}

