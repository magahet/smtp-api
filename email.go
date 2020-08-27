package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"smtp"
	"strconv"
	"strings"
	"text/template"
)

const (
	templatePath = "email.tmpl"
)

type message struct {
	from    string
	subject string
	text    string
	to      []string
}

type loginAuth struct {
	username, password string
}

// LoginAuth is used for smtp login auth
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func newMessage(from, subject, text string, to ...string) *message {
	return &message{
		from:    from,
		subject: subject,
		text:    text,
		to:      to,
	}
}

func (message *message) Body() []byte {
	body := ""
	body += fmt.Sprintf("From: %s\r\n", message.from)
	if len(message.to) > 0 {
		body += fmt.Sprintf("To: %s\r\n", strings.Join(message.to, ";"))
	}

	body += fmt.Sprintf("Subject: %s\r\n", message.subject)
	body += "\r\n" + message.text

	return []byte(body)
}

func sendEmail(event *Event) error {

	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	to := os.Getenv("SMTP_TO")
	log.Println(to)

	// Render message body with email template
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var msg bytes.Buffer
	if err := t.Execute(&msg, event); err != nil {
		return err
	}

	subject := fmt.Sprintf("[%s] %s", event.Type, event.Summary)
	message := newMessage(smtpUser, subject, msg.String(), to)
	log.Println(string(message.Body()))

	auth := LoginAuth(smtpUser, smtpPass)
	hostString := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	err = smtp.SendMail(hostString, auth, smtpUser, message.to, message.Body())
	if err != nil {
		return err
	}

	log.Println("Mail sent successfully")
	return nil
}
