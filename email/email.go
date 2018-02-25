package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"mta_app/config"
	"net/smtp"
)

type emailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
	SendTo      []string
}

func NewEmailUser(opts config.EmailUser) emailUser {
	return emailUser{
		Username:    opts.Username,
		Password:    opts.Password,
		EmailServer: opts.Server,
		Port:        opts.Port,
		SendTo:      opts.SendTo,
	}
}

func (user emailUser) SendEmail(mtaInfo config.MTAInfo) error {
	auth := smtp.PlainAuth("", user.Username, user.Password, user.EmailServer)

	email, err := user.createEmail(mtaInfo)
	if err != nil {
		return err
	}

	if err := smtp.SendMail("smtp.gmail.com:587", auth, "benraskin92@gmail.com", user.SendTo, email); err != nil {
		return err
	}
	log.Printf("sending email to %s", user.SendTo)

	return nil
}

type SmtpTemplateData struct {
	From string
	// To      string
	Subject string
	Body    string
}

const emailTemplate = `From: {{.From}}
Subject: {{.Subject}}

{{.Body}}

{{.From}}
`

func (user emailUser) createEmail(mtaInfo config.MTAInfo) ([]byte, error) {
	var doc bytes.Buffer
	var email []byte

	context := &SmtpTemplateData{
		From:    "Ben",
		Subject: "Bus is 2 stops away!",
		Body:    fmt.Sprintf("M11 is currently at %s", mtaInfo.StopCheck),
	}

	t := template.New("emailTemplate")
	t, err := t.Parse(emailTemplate)
	if err != nil {
		err = fmt.Errorf("error trying to parse mail template: %v", err)
		return email, err
	}
	if err = t.Execute(&doc, context); err != nil {
		err = fmt.Errorf("error trying to execute mail template: %v", err)
		return email, err
	}
	email = doc.Bytes()
	return email, nil
}
