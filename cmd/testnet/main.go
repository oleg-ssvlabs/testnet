package main

import (
	"context"
	"log/slog"

	"github.com/oleg-ssvlabs/testnet/internal/blockchain"
	"github.com/oleg-ssvlabs/testnet/internal/logger"
)

func main() {
	ctx := context.Background()
	logger.Initialize(slog.LevelInfo)

	slog.Info("starting blockchain service")
	response, err := blockchain.RunFromSDK(ctx)
	if err != nil {
		panic(err)
	}
	slog.
		With("response", response).
		Info("blockchain service started. Deploying SSV contracts")
}
