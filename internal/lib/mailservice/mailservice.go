package mailservice

import (
	"bytes"
	"go_echo/app/user/model/user"
	"go_echo/internal/config/env"
	"go_echo/internal/config/mailer"
	"go_echo/internal/util/helper"
	"time"

	"github.com/wneessen/go-mail"
	"golang.org/x/text/language"
)

func SendUserConfirmEmail(user *user.User, token string) {
	data := MakeMailTemplateData(map[string]any{
		"user":  *user,
		"token": token,
	})
	template := helper.GetTemplate("auth/register.gohtml")
	var html bytes.Buffer
	helper.MustErr(template.Execute(&html, data))
	SendEmail(
		user.Email,
		"Welcome to "+env.GetConfigInstance().AppName,
		html.String(),
	)
}

func MakeMailTemplateData(data map[string]any) interface{} {
	_, ok := data["Locale"]
	if !ok {
		data["Locale"] = language.English.String()
	}
	cfg := env.GetConfigInstance()
	data["AppStaticLink"] = cfg.AppURL + "/" + cfg.Static.URL + "/images/"
	data["AppLink"] = cfg.AppURL
	data["AppName"] = cfg.AppName
	data["Year"] = time.Now().UTC().Year()
	return helper.MakeStruct(data)
}

func SendEmail(
	to string,
	subject string,
	body string,
) {
	cfg := env.GetConfigInstance()
	m := mail.NewMsg()
	helper.MustErr(m.From(cfg.Mail.AddressFrom))
	helper.MustErr(m.To(to))

	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)

	mailClient := mailer.GetMailInstance()
	if err := mailClient.DialAndSend(m); err != nil {
		panic("failed to send mail: " + err.Error())
	}
}
