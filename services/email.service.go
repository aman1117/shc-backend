package services

import (
	"bytes"
	"html/template"
	"net/smtp"
	"os"
	"strings"
)

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

type EmailService struct {
	smtpServer *smtpServer
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpServer: &smtpServer{
			host: "smtp.gmail.com", port: "587",
		},
	}
}

func (es *EmailService) SendEmail(
	to []string,
	subject string,
	body string,
) error {
	from := "arl0817osho@gmail.com"
	password := os.Getenv("GOOGLE_APP_PASSWORD")
	auth := smtp.PlainAuth("shc-backend", from, password, es.smtpServer.host)

	// Create a new template
	tmpl := template.New("emailTemplate")

	// Parse the HTML email template
	templateString := `
		<html>
			<body>
				<h1>{{.Subject}}</h1>
				<p>{{.Body}}</p>
			</body>
		</html>
	`
	tmpl, err := tmpl.Parse(templateString)
	if err != nil {
		return err
	}

	// Prepare the data for the template
	data := struct {
		Subject string
		Body    string
	}{
		Subject: subject,
		Body:    body,
	}

	// Render the template into a buffer
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, data)
	if err != nil {
		return err
	}

	// Convert the buffer to a string
	htmlBody := buffer.String()

	// Set the email content type to HTML
	msg := "From: " + from + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		htmlBody

	// Send the email
	err = smtp.SendMail(es.smtpServer.Address(), auth, from, to, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
