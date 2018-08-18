package main

import (
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

func (s *Sender) Send(message *Message) error {
	conn, err := smtp.Dial(s.HostString())
	if err != nil {
		return err
	}

	// step 2: add all from and to
	if err = conn.Mail(message.from); err != nil {
		return err
	}
	for _, k := range message.to {
		if err = conn.Rcpt(k); err != nil {
			return err
		}
	}

	// Data
	w, err := conn.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message.Body()))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	conn.Quit()

	log.Println("Mail sent successfully")
	return nil
}
