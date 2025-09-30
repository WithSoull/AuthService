package tokens

import (
	"github.com/WithSoull/AuthService/internal/model"
)

type TokenGenerator interface {
	GenerateAccessToken(info model.UserInfo) (string, error)
	GenerateRefreshToken(info model.UserInfo) (string, error)
	VerifyAccessToken(tokenStr string) (*model.UserClaims, error)
	VerifyRefreshToken(tokenStr string) (*model.UserClaims, error)
}
