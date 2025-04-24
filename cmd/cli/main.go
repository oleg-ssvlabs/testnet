package main

import (
	"github.com/oleg-ssvlabs/testnet/internal/blockchain"
	"github.com/oleg-ssvlabs/testnet/internal/logger"
)

func main() {
	logger.Initialize()

	if err := blockchain.Run(); err != nil {
		panic(err)
	}
}
