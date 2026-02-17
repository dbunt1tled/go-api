package server

import (
	"context"
	"fmt"

	"github.com/dbunt1tled/go-api/internal/grpc/server/provider"
	"github.com/dbunt1tled/go-api/internal/modules/user_notification"
)

type ChannelProvider interface {
	Subscribe(ctx context.Context, channel string, userID int64) error
	Publish(ctx context.Context, channel string, userID int64, data []byte) (*[]byte, error)
}

type ChannelProviderResolver struct {
	providers map[string]*ChannelProvider
}

func NewChannelProviderResolver(
	userNotificationService *user_notification.Service,
) *ChannelProviderResolver {
	r := ChannelProviderResolver{
		providers: make(map[string]*ChannelProvider),
	}
	r.RegisterProvider(provider.ChannelUser, provider.NewUserProvider())
	r.RegisterProvider(provider.ChannelRead, provider.NewReadProvider(userNotificationService))

	return &r
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
