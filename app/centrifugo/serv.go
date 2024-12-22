package centrifugo

import (
	"context"
	"encoding/json"
	"go_echo/app/centrifugo/param"
	"go_echo/app/user/model/user"
	"go_echo/app/user/service"
	"go_echo/internal/config/logger"
	proxyproto "go_echo/internal/grpc"
	"go_echo/internal/util/jwt"
	"strconv"

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
	err = json.Unmarshal(request.GetData(), &req)
	if err != nil {
		log.ErrorContext(ctx, "Centrifugo Connect error unmarshal request",
			"request_data", request.GetData(),
			"error", err,
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
			"request_data", request.GetData(),
			"error", err,
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidAccessToken,
				Message: "invalid access token",
			},
		}, nil
	}
	u, err = service.UserRepository{}.ByID(int64(token["iss"].(float64))) //nolint:errcheck
	if err != nil || u.Status != user.Active {

		log.ErrorContext(ctx, "Centrifugo Connect error find user by id",
			"request_data", request.GetData(),
			"error", err,
			"user", u,
		)
		return &proxyproto.ConnectResponse{
			Error: &proxyproto.Error{
				Code:    ErrInvalidUser,
				Message: "invalid user",
			},
		}, nil
	}
	userId := strconv.FormatInt(u.ID, 10)
	return &proxyproto.ConnectResponse{
		Result: &proxyproto.ConnectResult{
			User:     userId,
			Data:     []byte(`{"user_id":` + userId + `"channel": ["user:#` + userId + `", "read:#` + userId + `"]}`),
			ExpireAt: int64(token["exp"].(float64)),
		},
	}, nil
}

func (s *Server) Refresh(ctx context.Context, request *proxyproto.RefreshRequest) (*proxyproto.RefreshResponse, error) {

	return &proxyproto.RefreshResponse{}, nil
}

func (s *Server) Subscribe(ctx context.Context, request *proxyproto.SubscribeRequest) (*proxyproto.SubscribeResponse, error) {

	return &proxyproto.SubscribeResponse{}, nil
}

func (s *Server) Publish(ctx context.Context, request *proxyproto.PublishRequest) (*proxyproto.PublishResponse, error) {

	// log.Println(string(request.Data))
	return &proxyproto.PublishResponse{}, nil
}

func (s *Server) RPC(ctx context.Context, request *proxyproto.RPCRequest) (*proxyproto.RPCResponse, error) {

	return &proxyproto.RPCResponse{}, nil
}

func (s *Server) SubRefresh(context.Context, *proxyproto.SubRefreshRequest) (*proxyproto.SubRefreshResponse, error) {

	return &proxyproto.SubRefreshResponse{}, nil
}

func (s *Server) SubscribeUnidirectional(*proxyproto.SubscribeRequest, grpc.ServerStreamingServer[proxyproto.StreamSubscribeResponse]) error {

	return nil
}

func (s *Server) SubscribeBidirectional(grpc.BidiStreamingServer[proxyproto.StreamSubscribeRequest, proxyproto.StreamSubscribeResponse]) error {

	return nil
}

func (s *Server) NotifyCacheEmpty(context.Context, *proxyproto.NotifyCacheEmptyRequest) (*proxyproto.NotifyCacheEmptyResponse, error) {

	return &proxyproto.NotifyCacheEmptyResponse{}, nil
}
func (s *Server) NotifyChannelState(context.Context, *proxyproto.NotifyChannelStateRequest) (*proxyproto.NotifyChannelStateResponse, error) {

	return &proxyproto.NotifyChannelStateResponse{}, nil
}

func (s *Server) mustEmbedServer() {
	panic("implement me")
}