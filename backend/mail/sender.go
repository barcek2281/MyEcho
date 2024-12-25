package mail

import "net/smtp"

type Sender struct {
	emailTo         string
	emailToPassword string
}

func NewSender(emailFrom, emailPassword string) *Sender {
	return &Sender{
		emailTo:         emailFrom,
		emailToPassword: emailPassword,
	}
}

func (send *Sender) SendToSupport(head, body, who string) error {
	auth := smtp.PlainAuth("hitler", send.emailTo, send.emailToPassword, "smtp.gmail.com")

	to := []string{send.emailTo}

	msg := "Subject: " + head +
		"\r\n\n" +
		body + "\r\n" + who

	err := smtp.SendMail("smtp.gmail.com:587", auth, "", to, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (send *Sender) SendToEveryPerson(head, body string, people []string) error {
	auth := smtp.PlainAuth("hitler", send.emailTo, send.emailToPassword, "smtp.gmail.com")
	msg := "Subject: " + head +
		"\r\n\n" +
		body + "\r\n"
	err := smtp.SendMail("smtp.gmail.com:587", auth, "", people, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
