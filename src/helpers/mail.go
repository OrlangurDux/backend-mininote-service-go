package helpers

import (
	"net/smtp"

	middlewares "orlangur.link/services/mini.note/handlers"
)

//Mail -> message sender
func Mail(to []string, subject string, message string) error {
	host := middlewares.DotEnvVariable("SMTP_HOST")
	login := middlewares.DotEnvVariable("SMTP_LOGIN")
	password := middlewares.DotEnvVariable("SMTP_PASSWORD")
	port := middlewares.DotEnvVariable("SMTP_PORT")

	body := "Subject:" + subject + "\r\n\r\n" + message
	bBody := []byte(body)
	auth := smtp.PlainAuth("", login, password, host)

	err := smtp.SendMail(host+":"+port, auth, login, to, bBody)
	if err != nil {
		return err
	}
	return nil
}
