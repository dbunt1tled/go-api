package centrifugo

import (
	"context"
	"go_echo/app/centrifugo/param"
	"go_echo/app/centrifugo/server"
	"go_echo/app/user/model/user"
	"go_echo/app/user/service"
	"go_echo/internal/config/logger"
	proxyproto "go_echo/internal/grpc"
	"go_echo/internal/util/helper"
	"go_echo/internal/util/jwt"
	"log/slog"
	"strconv"

	"github.com/bytedance/sonic"
	"google.golang.org/grpc"
)

type Server struct {
	proxyproto.UnimplementedCentrifugoProxyServer
}

func (s *Server) Connect(
	ctx context.Context,
	request *proxyproto.ConnectRequest,
) (*proxyproto.ConnectResponse, error) {
	var (
		req   param.ConnectParam
		err   error
		token map[string]interface{}
		u     *user.User
	)
	log := logger.GetLoggerInstance()
	err = sonic.ConfigFastest.Unmarshal(request.GetData(), &req)
	if err != nil {
		log.ErrorContext(ctx, "Centrifugo Connect error unmarshal request",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidConnectRequest,
				Message: "invalid request",
			},
		}, nil
	}
	token, err = jwt.JWToken{}.Decode(req.AccessToken, true)
	if err != nil {
		log.ErrorContext(ctx, "Centrifugo Connect error decode user token",
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
	u, err = service.UserRepository{}.ByID(ctx, int64(token["iss"].(float64))) //nolint:nolintlint,errcheck
	if err != nil || u.Status != user.Active {

		log.ErrorContext(ctx, "Centrifugo Connect error find user by id",
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
	userID := strconv.FormatInt(u.ID, 10)
	data := map[string]interface{}{
		"channels": []string{"user:#" + userID, "read:#" + userID},
		"user":     u.FirstName + " " + u.SecondName,
	}
	dataBytes, err := sonic.ConfigFastest.Marshal(data)
	if err != nil {
		log.ErrorContext(ctx, "Centrifugo Connect error marshal data",
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
			User:     userID,
			Data:     dataBytes,
			ExpireAt: int64(token["exp"].(float64)), //nolint:errcheck
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
		userID   int64
	)
	log := logger.GetLoggerInstance()
	providerResolver := server.GetChannelProviderResolver()
	channel := request.GetChannel()
	userID, err = strconv.ParseInt(request.GetUser(), 10, 64)
	if err != nil {
		log.ErrorContext(ctx, "Centrifugo Subscribe error parse user id",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.SubscribeResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidSubScribeRequest,
				Message: "invalid request",
			},
		}, nil
	}

	provider, err = providerResolver.Resolve(helper.SubStr(channel, ":"))
	if err != nil {
		log.ErrorContext(ctx, "Centrifugo Subscribe error resolve provider",
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
		log.ErrorContext(ctx, "Centrifugo Subscribe error subscribe channel",
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
		userID   int64
		dt       *[]byte
	)
	log := logger.GetLoggerInstance()
	providerResolver := server.GetChannelProviderResolver()
	data := request.GetData()
	channel := request.GetChannel()
	userID, err = strconv.ParseInt(request.GetUser(), 10, 64)
	if err != nil {
		log.ErrorContext(ctx, "Centrifugo Publish error parse user id",
			err,
			slog.String("request_data", string(request.GetData())),
		)
		return &proxyproto.PublishResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidPublishRequest,
				Message: "invalid request",
			},
		}, nil
	}

	provider, err = providerResolver.Resolve(helper.SubStr(channel, ":"))
	if err != nil {
		log.ErrorContext(ctx, "Centrifugo Publish error resolve provider",
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
		log.ErrorContext(ctx, "Centrifugo Publish error publish channel",
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
