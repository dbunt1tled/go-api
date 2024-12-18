package mailservice

import (
	"bytes"
	"encoding/json"
	"go_echo/app/user/model/user"
	"go_echo/internal/config/env"
	"go_echo/internal/config/mailer"
	"go_echo/internal/rmq"
	"go_echo/internal/util/helper"
	"time"

	"github.com/wneessen/go-mail"
	"golang.org/x/text/language"
)

const (
	ConfirmSubject = "confirm"
)

type MailJobMessage struct {
	UserID  int    `json:"userId"`
	Subject string `json:"subject"`
	Token   string `json:"token,omitempty"`
}

func SendUserConfirm(userID int64, token string) {
	job := MailJobMessage{
		UserID:  int(userID),
		Subject: ConfirmSubject,
		Token:   token,
	}
	rc := rmq.GetRMQInstance(rmq.MailExchange)
	rc.Publish(rmq.MailExchange, rmq.MailQueue, string(helper.Must(json.Marshal(&job))))
}

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
