package grpc

import (
	"context"

	"github.com/WithSoull/AuthService/internal/model"
)

type UserClient interface {
	ValidateCredentials(context.Context, string, string) (model.ValidateCredentialsResult, error)
}
