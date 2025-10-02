package gomail

import (
	"strconv"

	"go.k6.io/k6/js/modules"
	"gopkg.in/gomail.v2"
)

func init() {
	modules.Register("k6/x/gomail", new(GoMail))
}

type options struct {
	Subject     string   `js:"subject"`
	Message     string   `js:"message"`
	Attachments []string `js:"attachments"`
	Recipients  []string `js:"recipients"`
	Cc          []string `js:"cc"`
	Bcc         []string `js:"bcc"`
}

type GoMail struct{}

func (g *GoMail) SendMail(host string, port string, sender string, passwrd string, options options) error {
	msg := gomail.NewMessage()
	msg.SetHeader("Subject", options.Subject)
	msg.SetHeader("From", sender)
	msg.SetHeader("To", options.Recipients...)

	if len(options.Cc) > 0 {
		msg.SetHeader("Cc", options.Cc...)
	}

	if len(options.Bcc) > 0 {
		msg.SetHeader("Bcc", options.Bcc...)
	}

	if len(options.Attachments) > 0 {
		for _, f := range options.Attachments {
			msg.Attach(f)
		}
	}

	msg.SetBody("text/html", options.Message)

	numPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(host, numPort, sender, passwrd)

	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
