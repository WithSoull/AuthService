package domainerrors

import (
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
)

var (
	// Authentication errors (Unauthenticated)
	ErrInvalidRefreshToken    = sys.NewCommonError("invalid refresh token", codes.Unauthenticated)
	ErrInvalidEmailOrPassword = sys.NewCommonError("invalid email or password", codes.Unauthenticated)

	// Rate limiting errors (ResourceExhausted)
	ErrTooManyAttempts = sys.NewCommonError("too many attempts", codes.ResourceExhausted)
)
