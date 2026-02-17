package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/dbunt1tled/go-api/internal/grpc/server/provider/readhandlers"
	"github.com/dbunt1tled/go-api/internal/modules/user_notification"
	"github.com/dbunt1tled/go-api/pkg/f"
	"github.com/pkg/errors"
)

const ChannelRead = "read"

type ReadProvider struct {
	resolver *ReadChannelResolver
}

func NewReadProvider(
	userNotificationService *user_notification.Service,
) *ReadProvider {
	return &ReadProvider{
		resolver: NewReadChannelResolver(userNotificationService),
	}
}

type ReadChannelHandler interface {
	Handle(ctx context.Context, userID int64, data []byte) (*[]byte, error)
}

type ReadChannelResolver struct {
	handlers map[string]*ReadChannelHandler
}

func NewReadChannelResolver(
	userNotificationService *user_notification.Service,
) *ReadChannelResolver {
	r := ReadChannelResolver{
		handlers: make(map[string]*ReadChannelHandler),
	}
	r.RegisterHandler(ChannelUser, readhandlers.NewUserReadChannelHandler(userNotificationService))
	return &r
}

func (u *ReadProvider) Subscribe(ctx context.Context, channel string, userID int64) error {
	if channel != (ChannelRead + ":#" + strconv.FormatInt(userID, 10)) {
		return errors.New("invalid read channel")
	}
	return nil
}

func (u *ReadProvider) Publish(ctx context.Context, channel string, userID int64, data []byte) (*[]byte, error) {
	var (
		err error
		ch  string
		ok  bool
	)
	dt := make(map[string]interface{})
	err = sonic.ConfigFastest.Unmarshal(data, &dt)
	if err != nil {
		return nil, errors.Wrap(err, "invalid read channel data")
	}
	ch, ok = dt["channel"].(string)
	if !ok || ch == "" {
		return nil, errors.New("invalid read channel")
	}
	handler, err := u.resolver.Resolve(f.SubStr(ch, ":"))
	if err != nil {
		return nil, err
	}
	return (*handler).Handle(ctx, userID, data)
}

func (r *ReadChannelResolver) RegisterHandler(channelName string, handler ReadChannelHandler) {
	r.handlers[channelName] = &handler
}

func (r *ReadChannelResolver) Resolve(channelName string) (*ReadChannelHandler, error) {
	handler, exists := r.handlers[channelName]
	if !exists {
		return nil, fmt.Errorf("read channel handler for %s not found", channelName)
	}
	return handler, nil
}
