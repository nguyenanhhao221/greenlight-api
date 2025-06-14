package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"time"

	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	mailClient *mail.Client
	sender     string
}

func New(host, username, password, sender string) (*Mailer, error) {
	client, err := mail.NewClient(host, mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(username), mail.WithPassword(password), mail.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	return &Mailer{
		mailClient: client,
		sender:     sender,
	}, nil
}

func (m Mailer) Send(userEmail string, templateFile string, data map[string]any) error {
	message := mail.NewMsg()
	if err := message.From(m.sender); err != nil {
		log.Printf("failed to set From address: %s\n", err)
		return err
	}
	if err := message.To(userEmail); err != nil {
		log.Printf("failed to set To address: %s\n", err)
		return err
	}

	// Set up email content to be sent via embed file and template
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		log.Printf("failed to ParseFS: %s\n", err)
		return err
	}

	// Set the email subject
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", nil)
	if err != nil {
		log.Printf("failed to ExecuteTemplate for subject: %s\n", err)
		return err
	}
	message.Subject(subject.String())

	// Set the email plainBody from the tmpl
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		log.Printf("failed to ExecuteTemplate for plainBody: %s\n", err)
		return err
	}
	message.SetBodyString(mail.TypeTextPlain, plainBody.String())

	// Alternative html body
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		log.Printf("failed to ExecuteTemplate for htmlBody: %s\n", err)
		return err
	}
	message.AddAlternativeString("text/html", htmlBody.String())

	if err := m.mailClient.DialAndSend(message); err != nil {
		log.Printf("failed to send mail: %s\n", err)
		return err
	}
	return nil
}
