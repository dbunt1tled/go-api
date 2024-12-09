package auth

import (
	"go_echo/app/auth/model/token"
	"go_echo/app/user/model/user"
	"go_echo/internal/config/env"
	"go_echo/internal/util/jwt"
	"time"
)

const (
	AccessTokenSubject  = "access_token"
	RefreshTokenSubject = "refresh_token"
	ConfirmTokenSubject = "confirm_token"
)

func GetAuthTokens(user user.User) (*token.Tokens, error) {
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	return &token.Tokens{
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

func GenerateConfirmToken(user user.User) (string, error) {
	cfg := env.GetConfigInstance()
	data := map[string]interface{}{
		"iss": user.ID,
		"sub": ConfirmTokenSubject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(cfg.JWT.ConfirmLifeTime).Unix(),
	}
	token, err := jwt.JWToken{}.Encode(data)
	if err != nil {
		return "", err
	}

	return token, nil
}
