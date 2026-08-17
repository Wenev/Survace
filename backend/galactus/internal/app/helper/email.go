package helper

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/smtp"
	"os"
	"regexp"
)

func SendEmail(recipient, otpCode string) error {
	host := os.Getenv("EMAIL_HOST")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	port := os.Getenv("EMAIL_PORT")

	subject := "SurVace Verification Code"
	body := fmt.Sprintf(
		"Hello,\n\n\n"+
			"Here is your Verification Code: %s\n\n",
		otpCode)
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		username, recipient, subject, body)

	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%s", host, port)
	err := smtp.SendMail(addr, auth, username, []string{recipient}, []byte(message))
	if err != nil {
		log.Print(err)
		return status.Error(codes.Internal, "Internal Server Error")
	}
	return nil
}

func IsEmail(str string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	return regexp.MustCompile(regex).MatchString(str)
}

func SendSuccessEmail(recipient, username, email, password, typ string) error {
	host := os.Getenv("EMAIL_HOST")
	sender := os.Getenv("EMAIL_USERNAME")
	senderPassword := os.Getenv("EMAIL_PASSWORD")
	port := os.Getenv("EMAIL_PORT")

	homepageURL := "http://localhost:5173/"
	loginURL := "http://localhost:5173/login"

	var subject, body, link string

	switch typ {
	case "login":
		subject = "Login Successful - Welcome Back!"
		link = homepageURL
		body = fmt.Sprintf(
			"Hello %s,\n\n"+
				"You have successfully logged in.\n\n"+
				"Your credentials:\n"+
				"Email: %s\n"+
				"Username: %s\n"+
				"Password: %s\n\n"+
				"Access your homepage here: %s\n\n"+
				"Thank you for using our service!",
			username, email, username, password, link)
	case "register":
		subject = "Registration Successful - Welcome!"
		link = loginURL
		body = fmt.Sprintf(
			"Hello %s,\n\n"+
				"Your registration was successful!\n\n"+
				"Your credentials:\n"+
				"Email: %s\n"+
				"Username: %s\n"+
				"Password: %s\n\n"+
				"Login to your account here: %s\n\n"+
				"Thank you for joining us!",
			username, email, username, password, link)
	case "reset":
		subject = "Password Reset Successful"
		link = loginURL
		body = fmt.Sprintf(
			"Hello %s,\n\n"+
				"Your password has been reset successfully.\n\n"+
				"Your credentials:\n"+
				"Email: %s\n"+
				"Username: %s\n"+
				"Password: %s\n\n"+
				"Login to your account here: %s\n\n"+
				"If you did not request this change, please contact support immediately.",
			username, email, username, password, link)
	default:
		subject = "Account Notification"
		body = "Unknown operation type."
	}

	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		sender, recipient, subject, body)

	auth := smtp.PlainAuth("", sender, senderPassword, host)
	addr := fmt.Sprintf("%s:%s", host, port)
	err := smtp.SendMail(addr, auth, sender, []string{recipient}, []byte(message))
	if err != nil {
		log.Print(err)
		return status.Error(codes.Internal, "Internal Server Error")
	}
	return nil
}
