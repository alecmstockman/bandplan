package clients

import (
	"bandplan/src/helpers"
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

func SendTestVerificationEmail() error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	code, err := helpers.GenerateSessionToken()
	if err != nil {
		return err
	}
	codeHash := helpers.HashRegistrationCode(code)

	fmt.Println("code: ", code)
	fmt.Println("codeHash: ", codeHash)

	sent, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_FROM"),
		To:      []string{"bandplanapp@gmail.com"},
		Subject: "BandPlan email test",
		Html: fmt.Sprintf(`  
		<body>
			<h1>Click this button to verify your account</h1>

			<a href="https://bandplan-production.up.railway.app/email/verification/page?token=%v">
			Verify your account
			</a>
		</body>`, code),
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

//go:embed templates/email_verification.html
var emailVerificationTemplateFiles embed.FS

var emailVerificationTemplate = template.Must(
	template.ParseFS(emailVerificationTemplateFiles, "templates/email_verification.html"),
)

func AccountVericationEmail(data models.AccountVerificationPage) error {
	var body bytes.Buffer

	if err := emailVerificationTemplate.ExecuteTemplate(
		&body,
		"email_verification.html",
		data,
	); err != nil {
		return fmt.Errorf("render new-user signup email: %w", err)
	}

	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	sent, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_FROM"),
		To:      []string{data.Email},
		Subject: "BandPlan: Please verify your email",
		Html:    body.String(),
	})
	if err != nil {
		return fmt.Errorf("send new-user verification email: %w", err)
	}

	log.Printf("email sent: %s", sent.Id)
	return nil
}
