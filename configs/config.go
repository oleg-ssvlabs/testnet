package configs

import (
	_ "embed"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed config.yaml
	configData []byte

	App Config
)

type Config struct {
	WithObservability bool `yaml:"with-observability"`
	WithLocalnet      bool `yaml:"with-localnet"`
}

func init() {
	if err := yaml.Unmarshal(configData, &App); err != nil {
		panic(fmt.Sprintf("failed to parse config: %v", err))
	}

	slog.Info("config loaded", slog.Any("cfg", App))
}
