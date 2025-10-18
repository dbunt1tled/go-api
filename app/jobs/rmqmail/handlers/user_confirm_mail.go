package handlers

import (
	"context"
	"fmt"
	"github.com/dbunt1tled/go-api/app/user/model/user"
	"github.com/dbunt1tled/go-api/app/user/service"
	"github.com/dbunt1tled/go-api/internal/lib/mailservice"
	"github.com/dbunt1tled/go-api/internal/rmq"
	"github.com/dbunt1tled/go-api/internal/util/helper"

	"github.com/bytedance/sonic"
)

const (
	ConfirmSubject = "confirm"
)

type UserConfirmationEmail struct {
}

type MailUserConfirmationJobMessage struct {
	UserID  int    `json:"userId"`
	Subject string `json:"subject"`
	Token   string `json:"token"`
}

func (e UserConfirmationEmail) Handle(ctx context.Context, body []byte) error {
	var (
		job MailUserConfirmationJobMessage
		err error
		u   *user.User
	)

	if err = sonic.ConfigFastest.Unmarshal(body, &job); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err.Error())
	}
	u, err = service.UserRepository{}.ByID(ctx, int64(job.UserID))
	if err != nil {
		return fmt.Errorf("user: #%d not found. %s", job.UserID, err.Error())
	}
	mailservice.SendUserConfirmEmail(u, job.Token)
	return nil
}

func (e UserConfirmationEmail) Send(userID int64, token string) {
	job := MailUserConfirmationJobMessage{
		UserID:  int(userID),
		Subject: ConfirmSubject,
		Token:   token,
	}
	rmq.Publish(
		rmq.MailExchange,
		rmq.MailQueue,
		ConfirmSubject,
		string(helper.Must(sonic.ConfigFastest.Marshal(&job))),
	)
}
