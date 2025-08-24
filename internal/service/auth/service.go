package auth

import (
	"PVZ/internal/cache"
	"PVZ/internal/repository"
	"PVZ/internal/utils"
)

type AuthService struct {
	blackList cache.BlackList
	tokenUtil utils.TokenUtils
	userRepo  repository.UserRepository
}

func NewAuthService(blackList cache.BlackList, tokenUtil utils.TokenUtils, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		blackList: blackList,
		tokenUtil: tokenUtil,
		userRepo:  userRepo,
	}
}
