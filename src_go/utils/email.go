package utils

import (
	"errors"
	"fmt"
	"net/mail"

	"XustAutoSignIn/config"
)

func GetMessage(from, to mail.Address, subject, body string) string {
	// setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Mime-Version"] = "1.0"
	headers["Content-Type"] = "text/html;charset=UTF-8"
	// setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	return message
}

func SendEmail(to mail.Address, subject, body string) error {
	cfg := config.GetConfig()
	client := cfg.GetEmailClient()
	if client == nil {
		return errors.New("Failed to get SMTP client")
	}
	defer client.Close()
	err := client.Rcpt(to.Address)
	if err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = w.Write([]byte(GetMessage(cfg.GetFromEmail(), to, subject, body)))
	if err != nil {
		return err
	}
	return nil
}
