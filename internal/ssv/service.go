package ssv

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/moby/go-archive"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	dockerfilePath  = "internal/ssv"
	dockerImageName = "localssv/ssv-network"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Deploy(ctx context.Context, executionNodeURL, deployerPK string) error {
	docker, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	tar, err := archive.TarWithOptions(dockerfilePath, &archive.TarOptions{})
	if err != nil {
		return err
	}
	res, err := docker.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Tags:       []string{dockerImageName},
		Dockerfile: "Dockerfile",
	})
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		slog.Debug(scanner.Text())
	}

	slog.Info("docker image was built. Starting container")

	command := []string{
		"forge",
		"script",
		"script/DeployAll.s.sol:DeployAll",
		"--broadcast",
		"--rpc-url",
		"${ETH_RPC_URL}",
		"--private-key",
		"${PRIVATE_KEY}",
		"--legacy",
		"--silent"}

	containerCreateResp, err := docker.ContainerCreate(ctx,
		&container.Config{
			Image: dockerImageName,
			Env:   containerEnvVars(executionNodeURL, deployerPK),
			Cmd:   []string{"/bin/sh", "-c", strings.Join(command, " ")},
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		&v1.Platform{},
		"ssv-network")
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	tar, err = archive.TarWithOptions("internal/ssv/contract", &archive.TarOptions{
		IncludeFiles: []string{"RegisterOperators.sol"},
	})
	if err != nil {
		return err
	}
	if err := docker.CopyToContainer(ctx, containerCreateResp.ID, "/app/", tar, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}); err != nil {
		return fmt.Errorf("failed to copy file to container: %w", err)
	}

	slog.Info("docker container was created. Starting container")

	if err := docker.ContainerStart(ctx, containerCreateResp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	slog.Info("container started")

	return nil
}

func containerEnvVars(executionNodeURL, deployerPK string) []string {
	return []string{
		fmt.Sprintf("ETH_RPC_URL=%s", executionNodeURL),
		fmt.Sprintf("PRIVATE_KEY=%s", deployerPK),
		"MINIMUM_BLOCKS_BEFORE_LIQUIDATION=100800",
		"MINIMUM_LIQUIDATION_COLLATERAL=200000000",
		"OPERATOR_MAX_FEE_INCREASE=3",
		"DECLARE_OPERATOR_FEE_PERIOD=259200",
		"EXECUTE_OPERATOR_FEE_PERIOD=345600",
		"VALIDATORS_PER_OPERATOR_LIMIT=500",
		"OPERATOR_KEYS_FILE=/app/operator_keys.json",
	}
}
