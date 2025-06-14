package shared

import "github.com/docker/docker/api/types/container"

const NetworkName = "observability-net"

var (
	NetworkMode = container.NetworkMode(NetworkName)
	Labels      = map[string]string{"stack": "localnet-observability"}
)
