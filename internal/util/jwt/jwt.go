package jwt

import (
	"crypto/ecdsa"
	"encoding/base64"
	"go_echo/internal/config/env"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type JWToken struct{}

// Encode time.Now().Add(time.Second * time.Duration(exp)).Unix(),
func (JWToken) Encode(payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims(payload)
	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	privateKey, err := getPrivateKey()
	if err != nil {
		return "", err
	}
	return token.SignedString(privateKey)
}

func (JWToken) Decode(token string, checkExpire bool) (map[string]interface{}, error) {
	var claims jwt.MapClaims

	tokenData, err := jwt.ParseWithClaims(token, &claims, func(_ *jwt.Token) (interface{}, error) {
		return getPublicKey()
	})

	if !tokenData.Valid || err != nil {
		if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorExpired != 0 && !checkExpire { //nolint:errcheck
			return claims, nil
		}
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func getPrivateKey() (*ecdsa.PrivateKey, error) {
	cfg := env.GetConfigInstance()

	privateKeyBytes, err := base64.StdEncoding.DecodeString(cfg.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func getPublicKey() (*ecdsa.PublicKey, error) {
	cfg := env.GetConfigInstance()

	publicKeyBytes, err := base64.StdEncoding.DecodeString(cfg.JWT.PublicKey)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
