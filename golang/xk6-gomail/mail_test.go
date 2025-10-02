package gomail

import (
	"os"
	"testing"
)

func TestSendMain(t *testing.T) {
	gomail := GoMail{}

	mailOptions := options{
		Subject: "Test Send Email With Gomail",
		Message: "This test only testing sending email with gomail",
		Attachments: []string{
			"./test_file/test_attachment.txt",
			"./test_file/test_attachment.jpg",
			"./test_file/test_attachment.pdf",
		},
		Recipients: []string{
			os.Getenv("RECEIVER_EMAIL"),
		},
	}

	gomail.SendMail(
		"smtp.gmail.com",
		"587",
		os.Getenv("SENDER_EMAIL"),
		os.Getenv("SENDER_PASSWORD"),
		mailOptions,
	)

}
