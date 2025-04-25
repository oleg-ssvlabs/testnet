package ssv

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/moby/go-archive"
)

const (
	dockerfilePath = "internal/ssv"
	dockerImageTag = "localssv/ssv-network"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Deploy(ctx context.Context, el string, privateKey string) error {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	tar, err := archive.TarWithOptions(dockerfilePath, &archive.TarOptions{})
	if err != nil {
		return err
	}
	res, err := cli.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Tags:       []string{dockerImageTag},
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

	return nil
}
