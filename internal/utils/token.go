package utils

import (
	"PVZ/internal/cache"
	"PVZ/internal/model"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

type blackListToken struct {
	blackList  cache.BlackList
	secretKey  []byte
	expiration time.Duration
}

func NewTokenService(blackList cache.BlackList, secretKey string, expiration time.Duration) *blackListToken {
	return &blackListToken{
		blackList:  blackList,
		secretKey:  []byte(secretKey),
		expiration: expiration,
	}
}

func (s *blackListToken) GenerateToken(userType string) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
		},
		Role: userType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *blackListToken) ValidateToken(tokenString string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok || claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	// Проверка чёрного списка
	isBlacklisted, err := s.blackList.IsTokenBlacklisted(tokenString)
	if err != nil {
		zap.L().Error("Error checking blacklisted tokens:", zap.Error(err))
		return nil, errors.New("failed to check blacklisted tokens")
	}
	if isBlacklisted {
		return nil, errors.New("token is blacklisted")
	}

	return claims, nil
}
