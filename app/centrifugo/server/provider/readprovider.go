package provider

import (
	"encoding/json"
	"fmt"
	"go_echo/app/centrifugo/server/provider/readhandlers"
	"go_echo/internal/util/helper"
	"strconv"

	"github.com/pkg/errors"
)

const channelReadName = "read:#"

var instanceHandler *ReadChannelResolver

type ReadProvider struct {
}

type ReadChannelHandler interface {
	Handle(userID int64, data []byte) (*[]byte, error)
}

type ReadChannelResolver struct {
	handlers map[string]*ReadChannelHandler
}

func (u *ReadProvider) Subscribe(channel string, userID int64) error {
	if channel != (channelReadName + strconv.FormatInt(userID, 10)) {
		return errors.New("invalid read channel")
	}
	return nil
}

func (u *ReadProvider) Publish(channel string, userID int64, data []byte) (*[]byte, error) {
	var (
		err error
		ch  string
	)
	handlerResolver := GetReadChannelResolver()
	dt := make(map[string]interface{})
	err = json.Unmarshal(data, &dt)
	if err != nil {
		return nil, errors.Wrap(err, "invalid read channel data")
	}
	ch = dt["channel"].(string)
	if ch == "" {
		return nil, errors.New("invalid read channel")
	}
	handler, err := handlerResolver.Resolve(helper.SubStr(ch, ":"))
	if err != nil {
		return nil, err
	}
	return (*handler).Handle(userID, data)
}

func GetReadChannelResolver() *ReadChannelResolver {
	if instanceHandler == nil {
		r := ReadChannelResolver{
			handlers: make(map[string]*ReadChannelHandler),
		}
		r.RegisterHandler("user", &readhandlers.UserReadChannelHandler{})
		instanceHandler = &r
	}

	return instanceHandler
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
