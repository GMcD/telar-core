package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net"
	"net/smtp"

	"github.com/red-gold/telar-core/pkg/log"
)

var auth smtp.Auth

//Request struct
type Email struct {
	refEmail  string
	password  string
	smtpEmail string
}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func (email *Email) initEmail() {
	host, _, _ := net.SplitHostPort(email.smtpEmail)
	auth = smtp.PlainAuth("", email.refEmail, email.password, host)
}

func NewEmail(refEmail string, password string, smtpEmail string) *Email {
	return &Email{
		refEmail:  refEmail,
		password:  password,
		smtpEmail: smtpEmail,
	}
}

func NewEmailRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (email *Email) SendEmail(req *Request, tmplPath string, data interface{}) (bool, error) {
	log.Info("Initial email...")
	email.initEmail()

	log.Info("Start parsing html template...")
	err := req.parseTemplate(tmplPath, data)
	if err != nil {
		return false, fmt.Errorf("Error in parsing html template: %s", err.Error())
	}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + req.subject + "!\n"
	msg := []byte(subject + mime + "\n" + req.body)

	log.Info("Start sending email from %s to %s...", email.refEmail, req.to)
	errEmail := smtp.SendMail(email.smtpEmail, auth, email.refEmail, req.to, msg)
	if errEmail != nil {
		return false, fmt.Errorf("Error sending email: %s", errEmail.Error())
	}
	log.Info("Email sent from %s to %s...", email.refEmail, req.to)
	return true, nil
}

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	log.Info("HTML parsed data ", data)
	log.Info("HTML parsed body", r.body)
	return nil
}
