package mail

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/barcek2281/MyEcho/pkg/utils"
)

type Sender struct {
	emailTo         string
	emailToPassword string
}

var errLargeFile = errors.New("file too large XD")

func NewSender(emailFrom, emailPassword string) *Sender {
	return &Sender{
		emailTo:         emailFrom,
		emailToPassword: emailPassword,
	}
}

func (send *Sender) SendToSupport(subject, body, who, filename string, data *string) error {
	if send.sizeOfBase64(data) > 1*1024*1024*1024 {
		return errLargeFile
	}

	auth := smtp.PlainAuth("lol", send.emailTo, send.emailToPassword, "smtp.gmail.com")

	headers := "MIME-Version: 1.0\n" +
		"Content-Type: multipart/mixed; boundary=boundary\n\n"

	message := bytes.NewBuffer(nil)
	message.WriteString("Subject: " + subject + "\n")
	message.WriteString(headers)
	message.WriteString("--boundary\n")
	message.WriteString("Content-Type: text/plain; charset=\"utf-8\"\n\n")
	message.WriteString(body)
	message.WriteString("\n\n--boundary\n")
	message.WriteString(fmt.Sprintf("Content-Type: text/plain; name=\"%s\"\n", filename))
	message.WriteString("Content-Transfer-Encoding: base64\n")
	message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\n\n", filename))
	message.WriteString((*data)[utils.FindSymbol(data, ','):])
	message.WriteString("\n--boundary--")

	err := smtp.SendMail("smtp.gmail.com:587", auth, "", []string{send.emailTo}, message.Bytes())
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

func (send *Sender) SendToPerson(head, body string, person []string) error {
	auth := smtp.PlainAuth("hitler", send.emailTo, send.emailToPassword, "smtp.gmail.com")
	msg := "Subject: " + head +
		"\r\n\n" +
		body + "\r\n"
	err := smtp.SendMail("smtp.gmail.com:587", auth, "", person, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (send *Sender) sizeOfBase64(s *string) int {
	return 4 * (len(*s) / 3)
}
