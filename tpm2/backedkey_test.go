package tpm2

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	swtpm "github.com/foxboron/swtpm_test"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
)

func TPMForTesting() transport.TPMCloser {
	if os.Getenv("SKIP_TPM_TESTS") != "" {
		return nil
	}

	tempdir, err := os.MkdirTemp("", "swtpm_test")
	if err != nil {
		panic(err)
	}

	swtpm, err := swtpm.OpenSwtpm(tempdir)
	if err != nil {
		panic(err)
	}

	getCmd := tpm2.GetCapability{
		Capability: tpm2.TPMCapAlgs,
	}
	_, err = getCmd.Execute(swtpm)
	if err != nil {
		panic(err)
	}

	return swtpm
}

const TestKeyHandle = 0x81008231

func TestBackedP256Key(t *testing.T) {
	tpm := TPMForTesting()
	if tpm == nil {
		t.Skip("TPM tests are disabled")
	}
	defer tpm.Close()

	cfg := BackedP256KeyConfig{
		KeyHandle: TestKeyHandle,
		Password:  []byte("test\x00pass\x18word\xf3"),
	}

	_, err := LoadBackedP256Key(cfg, tpm)
	if err == nil {
		t.Fatalf("LoadBackedP256Key() before provisioning should have failed")
	}

	k, err := ProvisionBackedP256Key(cfg, tpm)
	if err != nil {
		t.Fatalf("ProvisionBackedP256Key() failed: %v", err)
	}

	der, err := x509.MarshalPKIXPublicKey(k.Public())
	if err != nil {
		t.Fatalf("Failed to marshal public key: %v", err)
	}
	bs := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	t.Logf("Public key:\n%s", string(bs))

	k2, err := LoadBackedP256Key(cfg, tpm)
	if err != nil {
		t.Fatalf("LoadBackedP256Key() failed: %v", err)
	}
	_ = k2
}
