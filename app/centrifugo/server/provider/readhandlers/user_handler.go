package readhandlers

import (
	"go_echo/app/usernotification/model/usernotification"
	"go_echo/app/usernotification/service"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

type UserReadMessage struct {
	ID      int64                   `json:"id"`
	UserID  int64                   `json:"userId,omitempty"`
	Channel string                  `json:"channel"`
	Data    []byte                  `json:"data,omitempty"`
	Status  usernotification.Status `json:"status"`
}

type UserReadChannelHandler struct {
}

type ReadChannelHandler interface {
	Handle(userID int64, data []byte) (*[]byte, error)
}

func (u *UserReadChannelHandler) Handle(userID int64, data []byte) (*[]byte, error) {
	var (
		dt  UserReadMessage
		err error
	)
	err = sonic.ConfigFastest.Unmarshal(data, &dt)
	if err != nil {
		return nil, errors.Wrap(err, "invalid read channel")
	}
	_, err = service.UserNotificationRepository{}.Update(
		dt.ID,
		service.UserNotificationParams{Status: &dt.Status},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error updating user notification")
	}
	return nil, nil //nolint:nilnil // response is not needed
}
