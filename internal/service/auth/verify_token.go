package auth

import (
	"context"
)

func (s *authService) ValidateToken(ctx context.Context, token string) error {
	_, err := s.tokenGenerator.VerifyAccessToken(token)
	if err != nil {
		return err
	}

	return nil
}
