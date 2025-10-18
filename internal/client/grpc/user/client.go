package user

import (
	"context"

	"github.com/WithSoull/AuthService/internal/client/grpc"
	"github.com/WithSoull/AuthService/internal/model"
	user_v1 "github.com/WithSoull/UserServer/pkg/user/v1"
)

type client struct {
	client user_v1.UserV1Client
}

func NewClient(grpcClient user_v1.UserV1Client) grpc.UserClient {
	return &client{
		client: grpcClient,
	}
}

func (c *client) ValidateCredentials(ctx context.Context, email, password string) (model.ValidateCredentialsResult, error) {
	resp, err := c.client.ValidateCredentials(ctx, &user_v1.ValidateCredentialsRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return model.ValidateCredentialsResult{}, err
	}

	return model.ValidateCredentialsResult{
		Valid:  resp.GetValid(),
		UserID: resp.GetUserId(),
	}, nil
}
