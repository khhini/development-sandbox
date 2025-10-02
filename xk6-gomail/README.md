# xk6-gomail
[k6](https://github.com/grafana/k6) extension to k6 extension to use [gomail](https://pkg.go.dev/gopkg.in/gomail.v2#section-readme) package for sending emails. Implemented using the [xk6](https://github.com/grafana/xk6) system.

## Build

```shell
git clone https://github.com/khhini/xk6-gomail
cd xk6-gomail
xk6 build --with github.com/khhini/xk6-gomail=.
```

## Example
An example of using the extension below. The options defined in the "mailOptions" variable are optional.

```javascript
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
        "smtp.gmail.com", // SMTP Server
		"587", // SMTP PORT
        "sender@gmail.com", // SENDER EMAIL
        "superdupersecretpassword", // SENDER PASSWORD
        mailOptions
    );
}
```