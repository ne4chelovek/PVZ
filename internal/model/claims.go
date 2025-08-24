package model

import jwt "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}
