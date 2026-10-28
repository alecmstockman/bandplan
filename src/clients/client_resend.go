package clients

import (
	"bandplan/src/helpers"
	"bandplan/src/models"
	"bytes"
	"embed"
	"fmt"
	"html"
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

	sent, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_FROM"),
		To:      []string{"bandplanapp@gmail.com"},
		Subject: "BandPlan email test",
		Html: fmt.Sprintf(`  

		<body style="margin: 0; padding: 0; background-color: #111214; color: #eeeeee; font-family: Arial, Helvetica, sans-serif;">
		<table role="presentation" width="100%%" cellspacing="0" cellpadding="0">
			<tr>
			<td align="center" style="padding: 40px 16px;">
				<table role="presentation" width="100%%" cellspacing="0" cellpadding="0"
				style="max-width: 480px; background-color: #1c1e23; border: 1px solid #292d31; border-radius: 12px;">
				<tr>
					<td style="padding: 32px;">
					<p style="margin: 0 0 28px; color: #eeeeee; font-size: 24px; font-weight: bold;">
						Band<span style="color: #C92A2A;">Plan</span>
					</p>

					<h1 style="margin: 0 0 16px; color: #eeeeee; font-size: 26px; line-height: 1.3;">
						You're almost in.
					</h1>

					<p style="margin: 0 0 28px; color: #bfc3c9; font-size: 16px; line-height: 1.6;">
						Verify your email to finish setting up your account
						and get your band organized.
					</p>

					<table role="presentation" cellspacing="0" cellpadding="0">
						<tr>
						<td bgcolor="#C92A2A" style="border-radius: 6px; text-align: center;">
							<a href="https://bandplan-production.up.railway.app/email/verification/page?token=%v"
							style="display: inline-block; padding: 15px 24px; border: 1px solid #C92A2A; border-radius: 6px; color: #ffffff; font-size: 16px; font-weight: bold; text-decoration: none;">
							Verify your account
							</a>
						</td>
						</tr>
					</table>

					<p style="margin: 28px 0 0; color: #9ca2ab; font-size: 13px; line-height: 1.6;">
						If you didn't create a BandPlan account, you can ignore this email.
					</p>
					</td>
				</tr>
				</table>

				<p style="margin: 20px 0 0; color: #9ca2ab; font-size: 12px;">
				The app for bands.
				</p>
			</td>
			</tr>
		</table>
		</body>`, code),
	})
	if err != nil {
		return err
	}

	log.Printf("email sent: %s", sent.Id)
	return nil
}

//go:embed templates/new-user-signup.html
var emailTemplateFiles embed.FS

var newUserSignupTemplate = template.Must(
	template.ParseFS(emailTemplateFiles, "templates/new-user-signup.html"),
)

func NewUserSignup(data models.NewUserSignupData) error {
	var body bytes.Buffer

	if err := newUserSignupTemplate.ExecuteTemplate(
		&body,
		"new-user-signup.html",
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

//go:embed templates/email-verification.html
var emailVerificationTemplateFiles embed.FS

var emailVerificationTemplate = template.Must(
	template.ParseFS(emailVerificationTemplateFiles, "templates/email-verification.html"),
)

func AccountVericationEmail(data models.AccountVerificationPage) error {
	fmt.Println("- AccountVericationEmail")
	var body bytes.Buffer

	if err := emailVerificationTemplate.ExecuteTemplate(
		&body,
		"email-verification.html",
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

func SendFeedbackEmail(user models.User, bandID, feedback string) error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	body := fmt.Sprintf(`
		<h1>User BandPlan User Feedback </h1>

		<p>User Name: %s</p>
		<p>User ID: %s</p>
		<p>User email: %s</p>
		<p>User bandID: %s</p>
		<p>Feedback: %s</p>
	`,
		html.EscapeString(user.DisplayName),
		html.EscapeString(user.UserID),
		html.EscapeString(user.Email),
		html.EscapeString(bandID),
		html.EscapeString(feedback),
	)

	sent, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_FROM"),
		To:      []string{"bandplanapp@gmail.com"},
		Subject: "BandPlan User Feedback",
		Html:    body,
	})
	if err != nil {
		return err
	}

	log.Printf("email sent: %s", sent.Id)
	return nil
}
