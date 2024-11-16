package security

import (
	"go_first/internal/config/env"
	"go_first/internal/lib/common/jwt"
	"go_first/internal/lib/common/model/security"
	"go_first/internal/lib/common/model/user"
	"time"
)

const (
	AccessTokenSubject  = "access_token"
	RefreshTokenSubject = "refresh_token"
)

func GetAuthTokens(user user.User) (*security.Tokens, error) {
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	return &security.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateAccessToken(user user.User) (string, error) {
	cfg := env.GetConfigInstance()
	data := map[string]interface{}{
		"iss": user.ID,
		"sub": AccessTokenSubject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(cfg.JWT.AccessLifeTime).Unix(),
	}
	token, err := jwt.JWToken{}.Encode(data)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateRefreshToken(user user.User) (string, error) {
	cfg := env.GetConfigInstance()
	data := map[string]interface{}{
		"iss": user.ID,
		"sub": RefreshTokenSubject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(cfg.JWT.RefreshLifeTime).Unix(),
	}
	token, err := jwt.JWToken{}.Encode(data)
	if err != nil {
		return "", err
	}

	return token, nil
}
