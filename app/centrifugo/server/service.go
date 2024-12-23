package server

import (
	"fmt"
	"go_echo/app/centrifugo/server/provider"
	"go_echo/internal/util/helper"
)

type ChannelProvider interface {
	Subscribe(channel string, userID int64) error
	Publish(channel string, userID int64, data []byte) (*[]byte, error)
}

type ChannelProviderResolver struct {
	providers map[string]*ChannelProvider
}

var instance *ChannelProviderResolver //nolint:gochecknoglobals // singleton

func GetChannelProviderResolver() *ChannelProviderResolver {
	if instance == nil {
		r := ChannelProviderResolver{
			providers: make(map[string]*ChannelProvider),
		}
		r.RegisterProvider("user", &provider.UserProvider{})
		instance = &r
	}

	return instance
}

func (r *ChannelProviderResolver) RegisterProvider(channelName string, provider ChannelProvider) {
	r.providers[channelName] = &provider
}

func (r *ChannelProviderResolver) Resolve(channelName string) (*ChannelProvider, error) {
	pdr, exists := r.providers[channelName]
	if !exists {
		return nil, fmt.Errorf("job provider for %s not found", channelName)
	}
	return pdr, nil
}

func (r *ChannelProviderResolver) GetChannelName(channel string) string {
	return helper.SubStr(channel, ":")
}
