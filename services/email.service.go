package services

import (
	"bytes"
	"html/template"
	"net/smtp"
	"os"
	"strings"
)

// why we are using this struct?✅
type smtpServer struct {
	host string
	port string
}

// why we are using this function? what is this Address?✅
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

type EmailService struct {
	smtpServer *smtpServer
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpServer: &smtpServer{
			// why we are using this host and port?✅
			host: "smtp.gmail.com", port: "587",
		},
	}
}

func (es *EmailService) SendEmail(
	// why this to is []string why not it is simply a string?✅
	to []string,
	subject string,
	body string,
) error {
	from := "ajaysharma.13122000@gmail.com"
	password := os.Getenv("GOOGLE_APP_PASSWORD")

	//what is this below line? what is plain auth? -> for authenticating with the smtp server✅
	auth := smtp.PlainAuth("shc-backend", from, password, es.smtpServer.host)

	//why we doing this? what is a template?
	tmpl := template.New("emailTemplate")

	//why we use backticks in go? -> for multi-line strings and it is called a raw string literal✅
	templateString := `
		<html>
			<body>
			<h2>{{.Subject}}</h2>
			<p>{{.Body}}</p>
			</body>
		</html>
	`
	tmpl, err := tmpl.Parse(templateString)
	if err != nil {
		return err
	}

	// Prepare the data for the template
	// explain the below struct✅
	data := struct {
		Subject string
		Body    string
	}{
		Subject: subject,
		Body:    body,
	}

	// Render the template into a buffer
	// what is buffer?
	buffer := new(bytes.Buffer)

	// what does tmpl.Execute do? ✅
	err = tmpl.Execute(buffer, data)
	if err != nil {
		return err
	}

	// Convert the buffer to a string
	htmlBody := buffer.String()

	// Set the email content type to HTML
	// what does "\r\n" do? -> carriage return and new line✅
	msg := "From: " + from + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		// what are we doing in the below line?
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		htmlBody

	// Send the email
	// what is SendMail? why we need to convert msg to []byte?✅
	err = smtp.SendMail(es.smtpServer.Address(), auth, from, to, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
