package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/bytedance/sonic"
	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/internal/lib/view"
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/mailer"
	"github.com/google/uuid"
)

const (
	UserConfirmationEmailJobAction = "confirm"
)

type UserConfirmationEmailJob struct {
	userService *user.Service
	mailService *mailer.Mailer
}

type MailUserConfirmationJobMessage struct {
	UserID uuid.UUID `json:"userId"`
	Token  string    `json:"token"`
}

func NewUserConfirmationEmailJob(
	userService *user.Service,
	mailService *mailer.Mailer,
) *UserConfirmationEmailJob {
	return &UserConfirmationEmailJob{
		userService: userService,
		mailService: mailService,
	}
}

func (e UserConfirmationEmailJob) Action() string {
	return UserConfirmationEmailJobAction
}

func (e UserConfirmationEmailJob) Handle(ctx context.Context, body []byte) error {
	var (
		job MailUserConfirmationJobMessage
		err error
		u   *user.User
	)

	if err = sonic.ConfigFastest.Unmarshal(body, &job); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err.Error())
	}
	u, err = e.userService.ByID(ctx, job.UserID)
	if err != nil {
		return fmt.Errorf("user: #%d not found. %s", job.UserID, err.Error())
	}
	if !u.IsPending() {
		log.Logger().Warn(
			"User is not in pending state",
			slog.String("user_id", u.ID.String()),
			slog.String("user_status", u.Status.String()),
			slog.String("msg", string(body)),
		)
		return nil
	}
	templ, err := view.GetTemplate("auth/register.gohtml")
	if err != nil {
		return err
	}
	if templ == nil {
		return errors.New("template not found")
	}
	var html bytes.Buffer
	err = templ.Execute(&html, view.MakeTemplateData(map[string]any{
		"User":  u,
		"Token": job.Token,
	}))
	if err != nil {
		return err
	}

	return e.mailService.SendEmail(
		ctx,
		u.Email,
		fmt.Sprintf("Welcome to %s", config.Get().Name),
		html.String(),
	)
}
