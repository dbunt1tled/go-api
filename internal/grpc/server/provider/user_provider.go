package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const ChannelUser = "user"

type UserProvider struct {
}

func NewUserProvider() *UserProvider {
	return &UserProvider{}
}

func (u *UserProvider) Subscribe(ctx context.Context, channel string, userID uuid.UUID) error {
	if channel != (ChannelUser + ":" + userID.String()) {
		return errors.New("invalid user channel")
	}
	return nil
}

func (u *UserProvider) Publish(ctx context.Context, channel string, userID uuid.UUID, data []byte) (*[]byte, error) {
	return nil, nil
}
