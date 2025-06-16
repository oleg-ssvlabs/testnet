package main

import (
	"context"
	"log/slog"

	"github.com/oleg-ssvlabs/testnet/configs"
	"github.com/oleg-ssvlabs/testnet/internal/localnet"
	"github.com/oleg-ssvlabs/testnet/internal/logger"
	"github.com/oleg-ssvlabs/testnet/internal/observability"
)

func main() {
	ctx := context.Background()
	logger.Initialize(slog.LevelInfo)

	if configs.App.WithLocalnet {
		slog.Info("starting localnet")
		err := localnet.Start(ctx)
		if err != nil {
			panic(err)
		}
	}

	slog.Info("localnet started.")

	if configs.App.WithObservability {
		slog.Info("Starting observability services")
		if err := observability.Start(ctx); err != nil {
			panic(err)
		}
	}

}
