package localnet

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	config "github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"

	_ "embed"
)

//go:embed params.yaml
var params []byte

const (
	parameterFile       = "params.yaml"
	enclaveName         = "localnet"
	kurtosisPackageName = "github.com/ssvlabs/ssv-mini"
)

func Start(ctx context.Context) error {
	kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		return errors.Join(err, errors.New("failed to create kurtosis context"))
	}

	enclaveCtx, err := kurtosisCtx.CreateEnclave(ctx, enclaveName)
	if err != nil {
		return errors.Join(err, errors.New("failed to create enclave"))
	}

	outputCh, cancel, err := enclaveCtx.RunStarlarkRemotePackage(
		ctx,
		kurtosisPackageName,
		config.NewRunStarlarkConfig(config.WithSerializedParams(string(params))))
	if err != nil {
		return errors.Join(err, errors.New("failed to run starlark package"))
	}
	defer cancel()

	var jsonResponse string
	for output := range outputCh {
		slog.Debug(output.String())

		ev := output.GetRunFinishedEvent()
		if ev != nil && ev.SerializedOutput != nil {
			jsonResponse = *ev.SerializedOutput
		}
	}

	var cfg any
	err = json.Unmarshal([]byte(jsonResponse), &cfg)
	if err != nil {
		return errors.Join(err, errors.New("failed to unmarshal json response"))
	}

	slog.Debug("response unmarshaled", slog.Any("response", cfg))

	return nil
}
