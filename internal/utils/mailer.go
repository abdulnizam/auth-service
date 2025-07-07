package utils

import (
	"fmt"
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

func sendMail(toEmail, subject, htmlBody string) error {
	mj := mailjet.NewMailjetClient(os.Getenv("MJ_API_KEY"), os.Getenv("MJ_SECRET_KEY"))

	messageInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: os.Getenv("EMAIL_FROM"),
				Name:  "Toggo App",
			},
			To: &mailjet.RecipientsV31{
				{
					Email: toEmail,
				},
			},
			Subject:  subject,
			TextPart: "Please view this email in HTML format.",
			HTMLPart: htmlBody,
		},
	}

	messages := mailjet.MessagesV31{Info: messageInfo}
	_, err := mj.SendMailV31(&messages)
	return err
}

// Send code only
func SendVerificationEmail(toEmail, token string) error {
	subject := "Verify your account"
	html := fmt.Sprintf(`
		<h3>Verify your account</h3>
		<p>Your verification code is: <b>%s</b></p>
	`, token)
	return sendMail(toEmail, subject, html)
}

// Send link to verify
func SendVerificationLinkEmail(toEmail, token string) error {
	verifyURL := fmt.Sprintf("http://localhost:3001/verify?token=%s&email=%s", token, toEmail)
	subject := "Verify your email"
	html := fmt.Sprintf(`
		<h2>Welcome!</h2>
		<p>Click the link below to verify your email:</p>
		<p><a href="%s" style="padding: 10px 20px; background-color: #0070f3; color: white; text-decoration: none;">Verify Email</a></p>
	`, verifyURL)
	return sendMail(toEmail, subject, html)
}
