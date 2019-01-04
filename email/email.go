package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

func sendMail(title string, body string, email []string) bool{
	auth := smtp.PlainAuth("", UserEmail, Mail_Password, Mail_Smtp_Host)
	to := email
	nickname := "秋风"
	user := UserEmail

	subject := title
	contentType := "Content-Type: text/plain; charset=UTF-8"
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	err := smtp.SendMail(Mail_Smtp_Host+Mail_Smtp_Port, auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
		return false
	}
	return true
}

