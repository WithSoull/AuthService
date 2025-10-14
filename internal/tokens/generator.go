package tokens

import (
	"context"

	"github.com/WithSoull/AuthService/internal/model"
)

type TokenGenerator interface {
	GenerateAccessToken(context.Context, model.UserInfo) (string, error)
	GenerateRefreshToken(context.Context, model.UserInfo) (string, error)
	VerifyAccessToken(context.Context, string) (*model.UserClaims, error)
	VerifyRefreshToken(context.Context, string) (*model.UserClaims, error)
}
