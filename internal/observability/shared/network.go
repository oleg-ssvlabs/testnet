package shared

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func EnsureNetwork(ctx context.Context, cli *client.Client) error {
	args := filters.NewArgs()
	args.Add("name", ObservabilityNetworkName)
	args.Add("name", LocalnetNetworkName)

	networks, err := cli.NetworkList(ctx, network.ListOptions{Filters: args})
	if err != nil {
		return errors.Join(err, errors.New("failed to list Docker network"))
	}

	var localnetAvailable, observabilityAvailable bool
	for _, network := range networks {
		if network.Name == LocalnetNetworkName {
			localnetAvailable = true
		}
		if network.Name == ObservabilityNetworkName {
			observabilityAvailable = true
		}
	}

	if !localnetAvailable {
		return fmt.Errorf("network: '%s' must be available before launching observability stack", LocalnetNetworkName)
	}

	if !observabilityAvailable {
		_, err = cli.NetworkCreate(ctx, ObservabilityNetworkName, network.CreateOptions{
			Driver: "bridge",
			Labels: Labels,
		})
		if err != nil {
			return errors.Join(err, errors.New("failed to create a network"))
		}
	}

	return nil
}
