package clients

import (
	"bandplan/src/models"
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/resend/resend-go/v4"
)

func SendTestEmail() error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	sent, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_FROM"),
		To:      []string{"bandplanapp@gmail.com"},
		Subject: "BandPlan email test",
		Html:    "<p>Email sending is working!</p>",
	})
	if err != nil {
		return err
	}

	log.Printf("email sent: %s", sent.Id)
	return nil
}

//go:embed templates/new_user_signup.html
var emailTemplateFiles embed.FS

var newUserSignupTemplate = template.Must(
	template.ParseFS(emailTemplateFiles, "templates/new_user_signup.html"),
)

func NewUserSignup(data models.NewUserSignupData) error {
	var body bytes.Buffer

	if err := newUserSignupTemplate.ExecuteTemplate(
		&body,
		"new_user_signup.html",
		data,
	); err != nil {
		return fmt.Errorf("render new-user signup email: %w", err)
	}

	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	sent, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_FROM"),
		To:      []string{"bandplanapp@gmail.com"},
		Subject: "BandPlan: New User Signup",
		Html:    body.String(),
	})
	if err != nil {
		return fmt.Errorf("send new-user signup email: %w", err)
	}

	log.Printf("email sent: %s", sent.Id)
	return nil
}
