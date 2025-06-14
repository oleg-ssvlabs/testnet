package observability

import (
	"context"
	"errors"
	"log/slog"

	"github.com/docker/docker/client"
	"github.com/oleg-ssvlabs/testnet/internal/observability/alloy"
	"github.com/oleg-ssvlabs/testnet/internal/observability/grafana"
	"github.com/oleg-ssvlabs/testnet/internal/observability/loki"
	"github.com/oleg-ssvlabs/testnet/internal/observability/shared"
)

func Start(ctx context.Context) error {
	slog.Info("instantiating Docker client")

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Join(err, errors.New("failed to instantiate Docker client"))
	}
	defer cli.Close()

	slog.With("network_name", shared.NetworkName).Info("creating new shared Docker network")
	networkID, err := shared.EnsureNetwork(ctx, cli)
	if err != nil {
		return errors.Join(err, errors.New("failed to create a Docker network"))
	}

	slog.With("network_id", networkID).Info("network was created")

	if err := grafana.Start(ctx, cli); err != nil {
		return errors.Join(err, errors.New("failed to start Grafana service"))
	}

	if err := loki.Start(ctx, cli); err != nil {
		return errors.Join(err, errors.New("failed to start Loki service"))
	}

	if err := alloy.Start(ctx, cli); err != nil {
		return errors.Join(err, errors.New("failed to start Alloy service"))
	}

	return nil
}
