package tpm2

import (
	"fmt"
	"os"

	swtpm "github.com/foxboron/swtpm_test"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
	"go.uber.org/zap"
)

const DefaultTPMPath = "/dev/tpmrm0"

type swtpmRmTmpDir struct {
	backing transport.TPMCloser
	tempdir string
}

func (s *swtpmRmTmpDir) Send(input []byte) ([]byte, error) {
	return s.backing.Send(input)
}

func (s *swtpmRmTmpDir) Close() error {
	if err := s.backing.Close(); err != nil {
		return err
	}
	return os.RemoveAll(s.tempdir)
}

func getSWTPM() (transport.TPMCloser, error) {
	s := zap.S()
	s.Error("!!! Using SWTPM instead of real TPM !!!")

	tempdir, err := os.MkdirTemp("", "h132_swtpm")
	if err != nil {
		return nil, fmt.Errorf("failed to create tempdir: %w", err)
	}

	swtpm, err := swtpm.OpenSwtpm(tempdir)
	if err != nil {
		return nil, fmt.Errorf("failed to open swtpm: %w", err)
	}

	return &swtpmRmTmpDir{
		backing: swtpm,
		tempdir: tempdir,
	}, nil
}

func GetTPM() (transport.TPMCloser, error) {
	s := zap.S()

	if os.Getenv("H132_USE_SWTPM") != "" {
		return getSWTPM()
	}
	tpmPath := os.Getenv("H132_TPM_PATH")
	if tpmPath == "" {
		tpmPath = DefaultTPMPath
	}
	if _, err := os.Stat(tpmPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("TPM device not found: %s", tpmPath)
	}

	s.Infof("Using TPM device: %s", tpmPath)
	t, err := transport.OpenTPM(tpmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open TPM device: %w", err)
	}

	getCmd := tpm2.GetCapability{
		Capability: tpm2.TPMCapAlgs,
	}
	_, err = getCmd.Execute(t)
	if err != nil {
		return nil, fmt.Errorf("failed to get TPM capabilities: %w", err)
	}

	return t, nil
}
