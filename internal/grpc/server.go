package centrifugo

import (
	"context"
	"log/slog"

	"github.com/bytedance/sonic"
	"github.com/dbunt1tled/go-api/internal/grpc/param"
	"github.com/dbunt1tled/go-api/internal/grpc/proxyproto"
	"github.com/dbunt1tled/go-api/internal/grpc/server"
	"github.com/dbunt1tled/go-api/internal/modules/auth"
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/internal/modules/user_notification"
	"github.com/dbunt1tled/go-api/pkg/f"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Server struct {
	proxyproto.UnimplementedCentrifugoProxyServer

	authService             *auth.Service
	userService             *user.Service
	userNotificationService *user_notification.Service
	channelProviderResolver *server.ChannelProviderResolver
}

func NewCentrifugoServer(
	authService *auth.Service,
	userService *user.Service,
	userNotificationService *user_notification.Service,
) *Server {
	return &Server{
		authService:             authService,
		userService:             userService,
		userNotificationService: userNotificationService,
		channelProviderResolver: server.NewChannelProviderResolver(
			userNotificationService,
		),
	}
}

func (s *Server) Connect(
	ctx context.Context,
	request *proxyproto.ConnectRequest,
) (*proxyproto.ConnectResponse, error) {
	var (
		req       param.ConnectParam
		err       error
		id        uuid.UUID
		token     map[string]interface{}
		u         *user.User
		dataBytes []byte
	)
	err = sonic.ConfigFastest.Unmarshal(request.GetData(), &req)
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Connect error unmarshal request",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidConnectRequest,
				Message: "invalid connect request",
			},
		}, nil
	}
	token, err = s.authService.DecodeAuthToken(req.AccessToken)
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Connect error decode user token",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidAccessToken,
				Message: "invalid access token",
			},
		}, nil
	}
	id, err = uuid.Parse(token["iss"].(string))
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Connect error token is invalid",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidAccessTokenUserId,
				Message: "invalid access token",
			},
		}, nil
	}

	u, err = s.userService.ByID(ctx, id)
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Connect error find user by id",
			err,
			slog.String("request_data", string(request.GetData())),
			slog.Any("user", u),
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidUser,
				Message: "invalid user",
			},
		}, nil
	}

	if u.Status != user.Active {
		log.Logger().WarnContext(ctx, "Centrifugo Connect error user is not active",
			slog.String("request_data", string(request.GetData())),
			slog.Any("user", u),
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrUserInactive,
				Message: "invalid user",
			},
		}, nil
	}

	data := map[string]interface{}{
		"channels": []string{"user:" + u.ID.String(), "read:" + u.ID.String()},
		"user":     u.FirstName + " " + u.SecondName,
	}
	dataBytes, err = sonic.ConfigFastest.Marshal(data)
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Connect error marshal data",
			err,
			slog.String("request_data", string(request.GetData())),
			slog.Any("user", u),
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidData,
				Message: "invalid data",
			},
		}, nil
	}
	return &proxyproto.ConnectResponse{
		Result: &proxyproto.ConnectResult{
			User:     u.ID.String(),
			Data:     dataBytes,
			ExpireAt: int64(token["exp"].(float64)),
		},
	}, nil
}

func (s *Server) Subscribe(
	ctx context.Context,
	request *proxyproto.SubscribeRequest,
) (*proxyproto.SubscribeResponse, error) {
	var (
		provider *server.ChannelProvider
		err      error
		userID   uuid.UUID
	)
	channel := request.GetChannel()
	userID, err = uuid.Parse(request.GetUser())
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Subscribe error parse user id",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.SubscribeResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidSubScribeRequest,
				Message: "invalid subscribe request",
			},
		}, nil
	}
	provider, err = s.channelProviderResolver.Resolve(f.SubStr(channel, ":"))
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Subscribe error resolve provider",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.SubscribeResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidSubscribeChannelProvider,
				Message: "invalid channel provider",
			},
		}, nil
	}
	err = (*provider).Subscribe(ctx, channel, userID)
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Subscribe error subscribe channel",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.SubscribeResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidSubscribeChannel,
				Message: "invalid channel",
			},
		}, nil
	}

	return &proxyproto.SubscribeResponse{}, nil
}

func (s *Server) Publish(
	ctx context.Context,
	request *proxyproto.PublishRequest,
) (*proxyproto.PublishResponse, error) {
	var (
		provider *server.ChannelProvider
		err      error
		userID   uuid.UUID
		dt       *[]byte
	)
	data := request.GetData()
	channel := request.GetChannel()
	userID, err = uuid.Parse(request.GetUser())
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Publish error parse user id",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.PublishResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidPublishRequest,
				Message: "invalid publish request",
			},
		}, nil
	}

	provider, err = s.channelProviderResolver.Resolve(f.SubStr(channel, ":"))
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Publish error resolve provider",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.PublishResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidPublishChannelProvider,
				Message: "invalid channel provider",
			},
		}, nil
	}
	dt, err = (*provider).Publish(ctx, channel, userID, data)
	if err != nil {
		log.Logger().ErrorContext(ctx, "Centrifugo Publish error publish channel",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.PublishResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidPublishChannelData,
				Message: "invalid channel",
			},
		}, nil
	}
	if dt == nil {
		return &proxyproto.PublishResponse{
			Result: &proxyproto.PublishResult{
				SkipHistory: true,
			},
		}, nil
	}
	return &proxyproto.PublishResponse{
		Result: &proxyproto.PublishResult{
			SkipHistory: true,
			Data:        *dt,
		},
	}, nil
}

func (s *Server) Refresh(
	context.Context,
	*proxyproto.RefreshRequest,
) (*proxyproto.RefreshResponse, error) {
	return &proxyproto.RefreshResponse{}, nil
}

func (s *Server) RPC(
	context.Context,
	*proxyproto.RPCRequest,
) (*proxyproto.RPCResponse, error) {
	return &proxyproto.RPCResponse{}, nil
}

func (s *Server) SubRefresh(
	context.Context,
	*proxyproto.SubRefreshRequest,
) (*proxyproto.SubRefreshResponse, error) {
	return &proxyproto.SubRefreshResponse{}, nil
}

func (s *Server) SubscribeUnidirectional(
	*proxyproto.SubscribeRequest,
	grpc.ServerStreamingServer[proxyproto.StreamSubscribeResponse],
) error {
	return nil
}

func (s *Server) SubscribeBidirectional(
	grpc.BidiStreamingServer[proxyproto.StreamSubscribeRequest, proxyproto.StreamSubscribeResponse],
) error {
	return nil
}

func (s *Server) NotifyCacheEmpty(
	context.Context,
	*proxyproto.NotifyCacheEmptyRequest,
) (*proxyproto.NotifyCacheEmptyResponse, error) {
	return &proxyproto.NotifyCacheEmptyResponse{}, nil
}
func (s *Server) NotifyChannelState(
	context.Context,
	*proxyproto.NotifyChannelStateRequest,
) (*proxyproto.NotifyChannelStateResponse, error) {
	return &proxyproto.NotifyChannelStateResponse{}, nil
}

func (s *Server) mustEmbedServer() {
	panic("mustEmbedServer implement me")
}
