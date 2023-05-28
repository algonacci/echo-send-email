package main

import (
	"log"
	"net/http"
	"net/smtp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.POST("/send-email", sendEmail)
	e.Logger.Fatal(e.Start(":8080"))
}

type EmailRequest struct {
	To      string `json:"to" form:"to" validate:"required,email"`
	Subject string `json:"subject" form:"subject" validate:"required"`
	Message string `json:"message" form:"message" validate:"required"`
}

type EmailResponse struct {
	Message string `json:"message"`
}

func sendEmail(c echo.Context) error {
	req := new(EmailRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to parse request body")
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := sendEmailSMTP(req.To, req.Subject, req.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to send email")
	}

	response := EmailResponse{Message: "Email sent successfully"}
	return c.JSON(http.StatusOK, response)

}

func sendEmailSMTP(to, subject, message string) error {
	from := "test@test.com"
	fromName := "Megalogic"
	password := "testpassword"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Compose the email message with HTML formatting and interactive styling
	msg := "From: " + fromName + " <" + from + ">\n" + // Include sender name
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		"<html><head><style>" +
		"body { font-family: Arial, sans-serif; background-color: #f5f5f5; }" +
		"h1 { color: #333; }" +
		"p { color: #666; }" +
		"a { color: #007bff; text-decoration: none; }" +
		"</style></head>" +
		"<body>" +
		"<div style=\"background-color: #ffffff; padding: 20px; margin: 20px; border-radius: 5px;\">" +
		"<h1 style=\"text-align: center;\">" + subject + "</h1>" +
		"<p>" + message + "</p>" +
		"<p style=\"text-align: center;\"><a href=\"https://example.com\" target=\"_blank\">Click here</a> for more information.</p>" +
		"</div>" +
		"</body></html>"

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
