package mailer

import (
	"go_echo/internal/config/env"
	"sync"

	"github.com/wneessen/go-mail"
)

var (
	mailInstance *mail.Client //nolint:gochecknoglobals // singleton
	m            sync.Once    //nolint:gochecknoglobals // singleton
)

func GetMailInstance() *mail.Client {
	m.Do(func() {
		var err error
		cfg := env.GetConfigInstance()
		mailInstance, err = mail.NewClient(
			cfg.Mail.Host,
			mail.WithPort(cfg.Mail.Port),
			mail.WithTLSPortPolicy(mail.TLSMandatory),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(cfg.Mail.Username),
			mail.WithPassword(cfg.Mail.Password),
		)
		if err != nil {
			panic("failed to create mail client: " + err.Error())
		}
	})
	return mailInstance
}

func Close() {
	mailInstance.Close()
}
