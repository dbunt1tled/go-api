package provider

import (
	"strconv"

	"github.com/pkg/errors"
)

const channelUserName = "user:#"

type UserProvider struct {
}

func (u *UserProvider) Subscribe(channel string, userID int64) error {
	if channel != (channelUserName + strconv.FormatInt(userID, 10)) {
		return errors.New("invalid user channel")
	}
	return nil
}

func (u *UserProvider) Publish(channel string, userID int64, data []byte) (*[]byte, error) {
	return nil, nil
}
