package main

import (
	"context"
	"log/slog"

	"github.com/oleg-ssvlabs/testnet/internal/blockchain"
	"github.com/oleg-ssvlabs/testnet/internal/logger"
	"github.com/oleg-ssvlabs/testnet/internal/ssv"
)

const withChain = false

func main() {
	ctx := context.Background()
	logger.Initialize(slog.LevelDebug)

	if withChain {
		slog.Info("starting blockchain service")
		response, err := blockchain.RunFromSDK(ctx)
		if err != nil {
			panic(err)
		}
		slog.
			With("response", response).
			Info("blockchain service started. Deploying SSV contracts")
	}

	slog.Info("deploying SSV contracts")

	ssvService := ssv.NewService()
	if err := ssvService.Deploy(ctx, "", ""); err != nil {
		panic(err)
	}
}
