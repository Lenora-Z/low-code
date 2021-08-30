package email

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendMail(title, body string, emails []string) error
}

type emailService struct {
	host     string
	port     int
	sender   string
	password string
}

func NewEmailService(sender, addr, pwd string, port uint64) EmailService {
	this := new(emailService)
	this.host = addr
	this.port = int(port)
	this.sender = sender
	this.password = pwd
	return this
}

func (srv *emailService) SendMail(title, body string, receiver []string) error {
	m := gomail.NewMessage()
	m.SetHeader(`From`, srv.sender)
	m.SetHeader(`To`, receiver...)
	m.SetHeader(`Subject`, title)
	m.SetBody(`text/html`, "<html><pre>"+body+"</pre></html>")

	err := gomail.NewDialer(srv.host, srv.port, srv.sender, srv.password).DialAndSend(m)
	if err != nil {
		logrus.Errorf("Send Email Fail, %s", err.Error())
		return err
	}
	return nil
}
