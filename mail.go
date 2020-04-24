package main

import (
	"bytes"
	"log"
	"net/smtp"
	"os"
)

// sendMail sends email about entry notification.
func sendMail(id, org string) {
	from := os.Getenv("SEND_FROM")
	to := os.Getenv("SEND_TO")
	if len(to) == 0 {
		return
	}
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	server := os.Getenv("SMTP_SERVER")
	port := os.Getenv("SMTP_PORT")
	body := bytes.NewBufferString("Subject: [GitHub/" + org + "] ID 招待依頼\r\n")
	body.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	body.WriteString("\r\n")
	body.WriteString("次のユーザーから GitHub " + org + " 組織への招待依頼が届きました: https://github.com/" + id + "\r\n")
	body.WriteString("\r\n")
	body.WriteString("メンバー管理はこちら: https://github.com/orgs/" + org + "/people\r\n")
	auth := smtp.PlainAuth("", user, password, server)
	if err := smtp.SendMail(server+":"+port, auth, from, []string{to}, body.Bytes()); err != nil {
		log.Printf("Failed to send mail, %v", err)
	}
}
