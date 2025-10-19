package centservice

import (
	"context"
	"strconv"

	"github.com/dbunt1tled/go-api/app/usernotification/model/usernotification"
	"github.com/dbunt1tled/go-api/app/usernotification/service"
	"github.com/dbunt1tled/go-api/internal/config/cntrfgclient"
	"github.com/dbunt1tled/go-api/internal/util/helper"

	"github.com/bytedance/sonic"
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
		ctx,
		service.UserNotificationParams{UserID: &message.UserID, Data: &dt},
	)
	if err != nil {
		return gocent.PublishResult{}, errors.Wrap(err, "Error save User Notification message")
	}
	d, err = sonic.ConfigFastest.Marshal(un)
	if err != nil {
		return gocent.PublishResult{}, errors.Wrap(err, "Error marshal User Notification message")
	}
	publish, err = cent.Publish(ctx, "user:#"+strconv.FormatInt(message.UserID, 10), d)
	if err != nil {
		return gocent.PublishResult{}, errors.Wrap(err, "Error publish User Notification message")
	}

	return publish, nil
}
