package access

import (
	"context"
	"errors"

	desc "github.com/WithSoull/AuthService/pkg/access/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *accessHandler) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errors.New("depricated")
}
