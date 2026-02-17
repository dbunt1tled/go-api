package provider

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
)

const ChannelUser = "user"

type UserProvider struct {
}

func NewUserProvider() *UserProvider {
	return &UserProvider{}
}

func (u *UserProvider) Subscribe(ctx context.Context, channel string, userID int64) error {
	if channel != (ChannelUser + ":#" + strconv.FormatInt(userID, 10)) {
		return errors.New("invalid user channel")
	}
	return nil
}

func (u *UserProvider) Publish(ctx context.Context, channel string, userID int64, data []byte) (*[]byte, error) {
	return nil, nil
}
