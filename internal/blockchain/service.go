package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"

	config "github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"

	_ "embed"
)

//go:embed params.yaml
var params []byte

const (
	parameterFile   = "params.yaml"
	ethereumPackage = "github.com/ethpandaops/ethereum-package"
)

func Run() error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(
		"bash",
		"-c",
		kurtosisCmd())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.Join(err, errors.New("failed to launch blockchain service"))
	}

	fmt.Print(stdout.String())
	fmt.Print(stderr.String())

	return nil
}

func kurtosisCmd() string {
	return fmt.Sprintf("kurtosis run --enclave testnet %s $(cat %s)", ethereumPackage, parameterFile)
}

type BlockchainResponse struct {
	ConsensusNodeURL string
	ExecutionNodeRPC string
	ExecutionNodeWS  string
}

func RunFromSDK(ctx context.Context) (BlockchainResponse, error) {
	var blockchainResponse BlockchainResponse
	kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		return blockchainResponse, errors.Join(err, errors.New("failed to create kurtosis context"))
	}

	enclaveCtx, err := kurtosisCtx.CreateEnclave(ctx, "testnet")
	if err != nil {
		return blockchainResponse, errors.Join(err, errors.New("failed to create enclave"))
	}

	outputCh, cancel, err := enclaveCtx.RunStarlarkRemotePackage(
		ctx,
		ethereumPackage,
		config.NewRunStarlarkConfig(config.WithSerializedParams(string(params))))
	if err != nil {
		return blockchainResponse, errors.Join(err, errors.New("failed to run starlark package"))
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

	var cfg ConfigResponse
	err = json.Unmarshal([]byte(jsonResponse), &cfg)
	if err != nil {
		return blockchainResponse, errors.Join(err, errors.New("failed to unmarshal json response"))
	}

	slog.Debug("response unmarshaled", slog.Any("response", cfg))

	return buildResponse(cfg), nil
}

func buildResponse(cfg ConfigResponse) BlockchainResponse {
	participant := cfg.AllParticipants[0]

	return BlockchainResponse{
		ConsensusNodeURL: participant.CLContext.BeaconHTTPURL,
		ExecutionNodeRPC: participant.ELContext.RpcHttpUrl,
		ExecutionNodeWS:  participant.ELContext.WsUrl,
	}
}
