package tpm2

import (
	"fmt"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
)

// This file manages the storage root key (SRK) for h132 use in the TPM.
// Note: We don't store SRK in the TPM NV. Per [1], recreating SRK on the fly is
//       often faster than loading it from NV.
// [1]: https://github.com/google/go-tpm/issues/335#issuecomment-1621874577

func LoadH132StorageRootKey(t transport.TPM) (*tpm2.AuthHandle, error) {
	cp := &tpm2.CreatePrimary{
		PrimaryHandle: tpm2.TPMRHOwner,
		InSensitive: tpm2.TPM2BSensitiveCreate{Sensitive: &tpm2.TPMSSensitiveCreate{
			UserAuth: tpm2.TPM2BAuth{Buffer: []byte{}},
		}},
		InPublic: tpm2.New2B(tpm2.ECCSRKTemplate),
	}

	resp, err := cp.Execute(t)
	if err != nil {
		return nil, fmt.Errorf("Failed to CreatePrimary(SRK): %v", err)
	}

	return &tpm2.AuthHandle{
		Handle: resp.ObjectHandle,
		Name:   resp.Name,
		Auth:   tpm2.PasswordAuth(nil),
	}, nil
}
