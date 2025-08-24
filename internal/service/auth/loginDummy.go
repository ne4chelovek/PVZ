package auth

import (
	"context"
	"fmt"
)

var validRoles = map[string]bool{
	"employee":  true,
	"moderator": true,
}

func (s *AuthService) LoginDummy(ctx context.Context, role string) (string, error) {
	if !validRoles[role] {
		return "", fmt.Errorf("invalid role: %s", role)
	}
	return s.tokenUtil.GenerateToken(role)
}
