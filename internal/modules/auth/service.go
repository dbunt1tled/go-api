package auth

import (
	"strings"
	"time"

	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/hasher"
	"github.com/labstack/echo/v5"
)

const BearerSchema = "Bearer "
const ErrorTokenMsg = "error token"

type Service struct {
	hasher *hasher.Hasher
}

func NewAuthService(hasher *hasher.Hasher) *Service {
	return &Service{
		hasher: hasher,
	}
}

func (s *Service) GeneratePasswordHash(password string) (string, error) {
	return s.hasher.HashArgon(password)
}

func (s *Service) ValidatePassword(password string, encodedHash string) (bool, error) {
	return s.hasher.CompareArgon(password, encodedHash)
}

func (s *Service) GenerateAuthTokens(user *user.User) (string, string, error) {
	var (
		err             error
		access, refresh string
	)
	access, err = s.hasher.EncodeJWT(map[string]interface{}{
		"iss": user.ID,
		"sub": hasher.AccessTokenSubject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(config.Get().Server.JWT.Expire.Access).Unix(),
	})
	if err != nil {
		return "", "", err
	}

	refresh, err = s.hasher.EncodeJWT(map[string]interface{}{
		"iss": user.ID,
		"sub": hasher.RefreshTokenSubject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(config.Get().Server.JWT.Expire.Refresh).Unix(),
	})
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) DecodeBearerToken(authorization string, opts ...hasher.DecodeOpt) (map[string]interface{}, error) {
	if !strings.HasPrefix(authorization, BearerSchema) {
		return nil, e.NewUnauthorizedError("Unauthorized", e.Err401TokenNotFoundError)
	}
	bearerToken := strings.TrimPrefix(authorization, BearerSchema)

	token, err := s.DecodeToken(bearerToken, opts...)

	if err != nil {
		return nil, err
	}

	return token, err
}

func (s *Service) DecodeToken(token string, opts ...hasher.DecodeOpt) (map[string]interface{}, error) {
	if token == "" {
		return nil, e.NewUnprocessableEntityError(ErrorTokenMsg, e.Err422TokenEmptyError)
	}
	t, err := s.hasher.DecodeJWT(token, opts...)
	if err != nil {
		return nil, e.NewUnprocessableEntityError(ErrorTokenMsg, e.Err422TokenError)
	}

	return t, err
}

func (s *Service) GenerateConfirmToken(user *user.User) (string, error) {
	return s.hasher.EncodeJWT(map[string]interface{}{
		"iss": user.ID,
		"sub": hasher.ConfirmTokenSubject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(config.Get().Server.JWT.Expire.Confirm).Unix(),
	})
}

func (s *Service) DecodeAuthToken(tokenAuth string) (map[string]interface{}, error) {
	return s.DecodeToken(
		tokenAuth,
		hasher.WithExpire(true),
		hasher.WithSubject(hasher.AccessTokenSubject),
	)
}

func (s *Service) DecodeRefreshToken(tokenAuth string) (map[string]interface{}, error) {
	return s.DecodeToken(
		tokenAuth,
		hasher.WithExpire(true),
		hasher.WithSubject(hasher.RefreshTokenSubject),
	)
}

func (s *Service) DecodeConfirmToken(tokenConfirm string) (map[string]interface{}, error) {
	return s.DecodeToken(tokenConfirm, hasher.WithSubject(hasher.ConfirmTokenSubject))
}

func (s *Service) TokenFromAuthHeader(c *echo.Context) (string, bool) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", true
	}
	authToken := strings.TrimSpace(strings.Split(authHeader, BearerSchema)[1])
	if authToken == "" {
		return "", true
	}
	return authToken, false
}

func (s *Service) TokenFromQueryParam(key string, c *echo.Context) (string, bool) {
	authToken := c.QueryParam(key)
	if authToken == "" {
		return "", true
	}
	return authToken, false
}
