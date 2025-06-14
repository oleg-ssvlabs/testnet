package shared

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func EnsureNetwork(ctx context.Context, cli *client.Client) (string, error) {
	args := filters.NewArgs()
	args.Add("name", NetworkName)

	networks, err := cli.NetworkList(ctx, network.ListOptions{Filters: args})
	if err != nil {
		return "", errors.Join(err, errors.New("failed to list Docker network"))
	}

	if len(networks) > 0 {
		return networks[0].ID, nil
	}

	resp, err := cli.NetworkCreate(ctx, NetworkName, network.CreateOptions{
		Driver: "bridge",
		Labels: Labels,
	})
	if err != nil {
		return "", errors.Join(err, errors.New("failed to create a network"))
	}

	return resp.ID, nil
}
