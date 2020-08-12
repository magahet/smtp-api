package main

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Message struct {
	from    string
	subject string
	text    string
	to      []string
}

type Sender struct {
	host string
	port int
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

// NewSender is used for smtp sending
func NewSender(host string, port int) *Sender {
	return &Sender{host, port}
}

func (s *Sender) HostString() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (mail *Sender) NewMessage(from, subject, text string, to ...string) *Message {
	return &Message{
		from:    from,
		subject: subject,
		text:    text,
		to:      to,
	}
}

func (message *Message) Body() string {
	body := ""
	body += fmt.Sprintf("From: %s\r\n", message.from)
	if len(message.to) > 0 {
		body += fmt.Sprintf("To: %s\r\n", strings.Join(message.to, ";"))
	}

	body += fmt.Sprintf("Subject: %s\r\n", message.subject)
	body += "\r\n" + message.text

	return body
}

func (s *Sender) Send(message *Message, user string, pass string) error {
	auth := LoginAuth(user, pass)
	msg := []byte(message.Body())
	err := smtp.SendMail(s.HostString(), auth, message.from, message.to, msg)
	if err != nil {
		return err
	}

	log.Println("Mail sent successfully")
	return nil
}
