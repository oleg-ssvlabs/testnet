package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

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
