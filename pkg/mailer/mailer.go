package mailer

import (
	"context"

	"github.com/wneessen/go-mail"
)

type Mailer struct {
	client    *mail.Client
	fromEmail string
}

func NewMailer(host string, port int, username string, password string, fromEmail string) *Mailer {
	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)
	if err != nil {
		panic("failed to create mail client: " + err.Error())
	}

	return &Mailer{
		client:    client,
		fromEmail: fromEmail,
	}
}

func (m *Mailer) Send(messages ...*mail.Msg) error {
	return m.client.DialAndSend(messages...)
}

func (m *Mailer) SendCtx(ctx context.Context, messages ...*mail.Msg) error {
	return m.client.DialAndSendWithContext(ctx, messages...)
}

func (m *Mailer) Close() error {
	return m.client.Close()
}


func (m *Mailer) SendEmail(
	c context.Context,
	to string,
	subject string,
	body string,
) error {
	e := mail.NewMsg()

	err := e.From(m.fromEmail)
	if err != nil {
		return err
	}
	err = e.To(to)
	if err != nil {
		return err
	}

	e.Subject(subject)
	e.SetBodyString(mail.TypeTextHTML, body)

	if err = m.SendCtx(c, e); err != nil {
		return err
	}

	return nil
}
