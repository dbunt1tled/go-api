package readhandlers

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/dbunt1tled/go-api/internal/modules/user_notification"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserReadMessage struct {
	ID      uuid.UUID                `json:"id"`
	UserID  uuid.UUID                `json:"userId,omitempty"`
	Channel string                   `json:"channel"`
	Data    []byte                   `json:"data,omitempty"`
	Status  user_notification.Status `json:"status"`
}

type UserReadChannelHandler struct {
	userNotificationService *user_notification.Service
}

func NewUserReadChannelHandler(
	userNotificationService *user_notification.Service,
) *UserReadChannelHandler {
	return &UserReadChannelHandler{
		userNotificationService: userNotificationService,
	}
}

type ReadChannelHandler interface {
	Handle(ctx context.Context, userID int64, data []byte) (*[]byte, error)
}

func (u *UserReadChannelHandler) Handle(ctx context.Context, userID int64, data []byte) (*[]byte, error) {
	var (
		dt  UserReadMessage
		un  *user_notification.UserNotification
		err error
	)
	err = sonic.ConfigFastest.Unmarshal(data, &dt)
	if err != nil {
		return nil, errors.Wrap(err, "invalid read channel")
	}
	un, err = u.userNotificationService.ByID(ctx, dt.ID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting user notification")
	}
	un.Status = dt.Status
	_, err = u.userNotificationService.Update(ctx, un)
	if err != nil {
		return nil, errors.Wrap(err, "error updating user notification")
	}
	return nil, nil
}
