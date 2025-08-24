package utils

import "PVZ/internal/model"

type TokenUtils interface {
	GenerateToken(userType string) (string, error)
	ValidateToken(tokenString string) (*model.UserClaims, error)
}
