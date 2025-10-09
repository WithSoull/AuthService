package validatorInterceptor

import (
	"context"

	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func ValidateNotEmptyEmailAndPassword(email, password string) validate.Condition {
	return func(ctx context.Context) error {
		if email == "" || password == "" {
			return validate.NewValidationErrors("empty credentials")
		}

		return nil
	}
}
