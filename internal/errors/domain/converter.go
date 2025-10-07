package domainerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Return true if log is needed, and error ofc
func ToGRPCStatus(err error) (bool, error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return false, status.Error(codes.Unauthenticated, "Invalid credentials")
	case errors.Is(err, ErrTooManyAttempts):
		return false, status.Error(codes.ResourceExhausted, "Too many attempts. Try again later.")
	default:
		return true, status.Error(codes.Internal, "unknown internal error")
	}
}
