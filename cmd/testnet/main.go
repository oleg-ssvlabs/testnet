package main

import (
	"context"
	"log/slog"

	"github.com/oleg-ssvlabs/testnet/internal/logger"
	"github.com/oleg-ssvlabs/testnet/internal/network"
	"github.com/oleg-ssvlabs/testnet/internal/observability"
)

func main() {
	ctx := context.Background()
	logger.Initialize(slog.LevelDebug)

	slog.Info("starting network")
	err := network.Start(ctx)
	if err != nil {
		panic(err)
	}

	slog.Info("network service started. Starting observability services")

	if err := observability.Start(ctx); err != nil {
		panic(err)
	}
}
