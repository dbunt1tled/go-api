package centrifugo

const (
	ErrInvalidConnectRequest = 1 + iota
	ErrInvalidAccessToken
	ErrInvalidAccessTokenUserId
	ErrInvalidUser
	ErrUserInactive
	ErrInvalidData
	ErrInvalidSubScribeRequest
	ErrInvalidSubscribeChannelProvider
	ErrInvalidSubscribeChannel
	ErrInvalidPublishRequest
	ErrInvalidPublishChannelProvider
	ErrInvalidPublishChannelData
)
