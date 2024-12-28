package centservice

import (
	"context"
	"encoding/json"
	"go_echo/app/usernotification/model/usernotification"
	"go_echo/app/usernotification/service"
	"go_echo/internal/config/cntrfgclient"
	"go_echo/internal/util/helper"
	"strconv"

	"github.com/centrifugal/gocent/v3"
	"github.com/pkg/errors"
)

type UserNotification struct {
	UserID  int64
	Message string
}

func SendUserNotification(message UserNotification) (gocent.PublishResult, error) {
	var (
		un      *usernotification.UserNotification
		d       []byte
		publish gocent.PublishResult
	)
	cent := cntrfgclient.GetInstance()
	ctx := context.Background()

	dt, err := helper.StructToJSONField(message)
	if err != nil {
		return gocent.PublishResult{}, errors.Wrap(err, "Error convert User Notification message to map")
	}
	un, err = service.UserNotificationRepository{}.Create(
		service.UserNotificationParams{UserID: &message.UserID, Data: &dt},
	)
	if err != nil {
		return gocent.PublishResult{}, errors.Wrap(err, "Error save User Notification message")
	}
	d, err = json.Marshal(un)
	if err != nil {
		return gocent.PublishResult{}, errors.Wrap(err, "Error marshal User Notification message")
	}
	publish, err = cent.Publish(ctx, "user:#"+strconv.FormatInt(message.UserID, 10), d)
	if err != nil {
		return gocent.PublishResult{}, errors.Wrap(err, "Error publish User Notification message")
	}

	return publish, nil
}
