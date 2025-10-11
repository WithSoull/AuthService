package main

import (
	"context"
	"fmt"

	"github.com/WithSoull/AuthService/internal/app"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appCtx := context.Background()

	a, err := app.NewApp(appCtx)
	if err != nil {
		panic(fmt.Sprintf("failed to init app: %v", err))
	}

	if err := a.Run(); err != nil {
		logger.Fatal(appCtx, "failed to run app", zap.Error(err))
	}
}
