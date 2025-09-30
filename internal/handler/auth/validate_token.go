package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *authHandler) ValidateToken(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header not provided")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")

	err := h.service.ValidateToken(ctx, token)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Unauthenticated, "invalid token")
	}
	return &emptypb.Empty{}, nil
}
