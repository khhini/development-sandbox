import gomail from 'k6/x/gomail';

export default function () {
    const mailOptions = { 
        subject: "Example Test Send Email With Gomail",
        message: "Example This test only testing sending email with gomail",
        attachments: [
			"./test_file/test_attachment.txt",
			"./test_file/test_attachment.jpg",
			"./test_file/test_attachment.pdf",
        ],
		recipients: [
			"sojive2197@evimzo.com",
        ],
    }

    gomail.sendMail(
        "smtp.gmail.com",
		"587",
        "sender@gmail.com",
        "password",
        mailOptions
    );
}