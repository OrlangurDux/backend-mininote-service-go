package helpers

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	middlewares "orlangur.link/services/mini.note/handlers"
	"orlangur.link/services/mini.note/models"
)

//go:embed templates/*/*.html
var templateFS embed.FS

func RenderContact(locale Locale, data models.ContactData) (models.RenderTemplate, error) {
	path := fmt.Sprintf("templates/%s/contact.html", locale)
	tmpl, err := template.ParseFS(templateFS, path)
	if err != nil {
		return models.RenderTemplate{}, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return models.RenderTemplate{}, err
	}
	subject := "Request from %s"
	switch locale {
	case LocaleRu:
		subject = "Запрос от %s"
	}
	subject = fmt.Sprintf(subject, middlewares.DotEnvVariable("NAME", "MiniNote"))
	return models.RenderTemplate{
		Subject: subject,
		Body:    buf.String(),
	}, nil
}

func RecoveryPasswordData(locale Locale, data models.RecoveryPasswordData) (models.RenderTemplate, error) {
	path := fmt.Sprintf("templates/%s/recovery_password.html", locale)
	tmpl, err := template.ParseFS(templateFS, path)
	if err != nil {
		return models.RenderTemplate{}, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return models.RenderTemplate{}, err
	}
	subject := "Password recovery on %s"
	switch locale {
	case LocaleRu:
		subject = "Восстановление пароля на %s"
	}
	subject = fmt.Sprintf(subject, middlewares.DotEnvVariable("NAME", "MiniNote"))
	return models.RenderTemplate{
		Subject: subject,
		Body:    buf.String(),
	}, nil
}
