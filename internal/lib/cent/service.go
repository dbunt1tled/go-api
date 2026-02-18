package cent

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/centrifugal/gocent/v3"
	"github.com/dbunt1tled/go-api/pkg/e"
)

type CentrifugoService struct {
	client *gocent.Client
}

func NewCentrifugoService(apiKey string, apiURL string) *CentrifugoService {
	return &CentrifugoService{
		client: gocent.New(gocent.Config{
			Addr: apiURL + "/api",
			Key:  apiKey,
		}),
	}
}

func (s *CentrifugoService) SendUserMessage(
	ctx context.Context,
	message *UserMessage,
) (*gocent.PublishResult, error) {
	var (
		d       []byte
		err     error
		publish gocent.PublishResult
	)
	d, err = sonic.ConfigFastest.Marshal(message)
	if err != nil {
		return nil, e.NewUnprocessableEntityErrorWrap(
			"Error marshal User Notification message",
			e.Err422UserNotificationMarshalError,
			err,
		)
	}
	publish, err = s.client.Publish(ctx, "user:"+message.UserID.String(), d)
	if err != nil {
		return nil, e.NewUnprocessableEntityErrorWrap(
			"Error publish User Notification message",
			e.Err422UserNotificationPublishError,
			err,
		)
	}
	return &publish, nil
}
