package webauthnwrappedtpm

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
	"go.uber.org/zap"
	"golang.org/x/crypto/hkdf"

	"github.com/IPA-CyberLab/h132/cmd/h132/common"
	"github.com/IPA-CyberLab/h132/pb"
	h132_tpm2 "github.com/IPA-CyberLab/h132/tpm2"
	"github.com/IPA-CyberLab/h132/webauthn"
)

var (
	CredAlreadySetErr    = errors.New("credential is already set")
	CredUnavailableErr   = errors.New("credential is not available. You must complete registration first.")
	NotYetProvisionedErr = errors.New("key is not yet provisioned onto the TPM")
)

type WebAuthnWrappedTPMKey struct {
	ReflectorUrlStr string
	TpmKeyHandle    tpm2.TPMHandle
	PrfSalt         []byte
	HkdfSalt        []byte

	cred *webauthn.Credential
	pub  *ecdsa.PublicKey
}

func NewUnprovisioned(tpm transport.TPM, reflectorUrlStr string, tpmKeyHandle tpm2.TPMHandle) (*WebAuthnWrappedTPMKey, error) {
	k := &WebAuthnWrappedTPMKey{
		ReflectorUrlStr: reflectorUrlStr,
		TpmKeyHandle:    tpmKeyHandle,
		PrfSalt:         make([]byte, 32),
		HkdfSalt:        make([]byte, 32),
	}

	if _, err := io.ReadFull(rand.Reader, k.PrfSalt); err != nil {
		zap.S().Panicf("Failed to generate salt: %v", err)
	}
	if _, err := io.ReadFull(rand.Reader, k.HkdfSalt); err != nil {
		zap.S().Panicf("Failed to generate salt: %v", err)
	}

	// Check if the key is already provisioned on the TPM
	_, _, err := h132_tpm2.ReadPublic(tpm, tpmKeyHandle)
	if err == nil {
		return nil, fmt.Errorf("TPM key already exists at %x", tpmKeyHandle)
	}
	if !errors.Is(err, h132_tpm2.HandleErr) {
		return nil, fmt.Errorf("unexpected TPM error: %w", err)
	}

	return k, nil
}

func webauthnUserName(lwsName, keyName string) string {
	return fmt.Sprintf("%s@%s", keyName, lwsName)
}

func (wk *WebAuthnWrappedTPMKey) StartRegistration(lwsName, keyName string) (*webauthn.WebAuthnRegistrationSession, error) {
	if wk.cred != nil {
		return nil, CredAlreadySetErr
	}

	return webauthn.StartRegistration(
		wk.ReflectorUrlStr,
		webauthnUserName(lwsName, keyName))
}

func (wk *WebAuthnWrappedTPMKey) CompleteRegistration(sess *webauthn.WebAuthnRegistrationSession, zb64resp []byte) error {
	cred, err := sess.Complete(zb64resp)
	if err != nil {
		return fmt.Errorf("failed to complete registration: %w", err)
	}
	if cred == nil {
		return common.ErrAbort{}
	}

	if wk.cred != nil {
		return CredAlreadySetErr
	}
	wk.cred = cred
	return nil
}

func (wk *WebAuthnWrappedTPMKey) StartGetPRFSecretSession() (*webauthn.GetPRFSecretSession, error) {
	return webauthn.StartGetPRFSecretSession(wk.cred, wk.PrfSalt, rand.Reader)
}

func derivePassword(prf, salt []byte, tpmKeyHandle tpm2.TPMHandle) []byte {
	hkdfr := hkdf.New(sha256.New, prf, salt, []byte(fmt.Sprintf("password %x", tpmKeyHandle)))
	tpmPassword := make([]byte, 32)
	if _, err := hkdfr.Read(tpmPassword); err != nil {
		zap.S().Panicf("Failed to generate tpmPassphrase: %v", err)
	}
	return tpmPassword
}

func (wk *WebAuthnWrappedTPMKey) Provision(tpm transport.TPM, prf []byte) error {
	tpmPassword := derivePassword(prf, wk.HkdfSalt, wk.TpmKeyHandle)

	cfg := h132_tpm2.BackedP256KeyConfig{
		KeyHandle: wk.TpmKeyHandle,
		Password:  tpmPassword,
	}
	bk, err := h132_tpm2.ProvisionBackedP256Key(cfg, tpm)
	if err != nil {
		return fmt.Errorf("failed to provision backed key: %w", err)
	}

	wk.pub = bk.Public()
	return nil
}

func GetProvisioned(tpm transport.TPM, lwsName, keyName string, p *pb.WebAuthnWrappedTPMKey) (*WebAuthnWrappedTPMKey, error) {
	waUserName := webauthnUserName(lwsName, keyName)

	wk := &WebAuthnWrappedTPMKey{
		ReflectorUrlStr: p.ReflectorUrl,
		TpmKeyHandle:    tpm2.TPMHandle(p.TpmKeyHandle),
		PrfSalt:         p.PrfSalt,
		HkdfSalt:        p.HkdfSalt,
		cred: &webauthn.Credential{
			UserName:        waUserName,
			ReflectorUrlStr: p.ReflectorUrl,
			WacJson:         p.WebauthnCredentialJson,
		},
	}

	_, pub, err := h132_tpm2.ReadPublic(tpm, wk.TpmKeyHandle)
	if errors.Is(err, h132_tpm2.HandleErr) {
		return nil, fmt.Errorf("key handle error. It is likely key doesn't exist at the specified handle: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("unexpected TPM error: %w", err)
	}
	wk.pub = pub

	return wk, nil
}

func (wk *WebAuthnWrappedTPMKey) Pub() (*ecdsa.PublicKey, error) {
	if wk.pub == nil {
		return nil, NotYetProvisionedErr
	}

	return wk.pub, nil
}

func (wk *WebAuthnWrappedTPMKey) SetImplProto(p *pb.KeyImpl) error {
	if wk.cred == nil {
		return CredUnavailableErr
	}

	p.Impl = &pb.KeyImpl_WebauthnWrappedTpm{
		WebauthnWrappedTpm: &pb.WebAuthnWrappedTPMKey{
			ReflectorUrl:           wk.ReflectorUrlStr,
			TpmKeyHandle:           uint32(wk.TpmKeyHandle),
			PrfSalt:                wk.PrfSalt,
			HkdfSalt:               wk.HkdfSalt,
			WebauthnCredentialJson: wk.cred.WacJson,
		},
	}
	return nil
}

func (wk *WebAuthnWrappedTPMKey) Unwrap(tpm transport.TPM, prf []byte) (*h132_tpm2.BackedP256Key, error) {
	tpmPassword := derivePassword(prf, wk.HkdfSalt, wk.TpmKeyHandle)
	cfg := h132_tpm2.BackedP256KeyConfig{
		KeyHandle: wk.TpmKeyHandle,
		Password:  tpmPassword,
	}
	return h132_tpm2.LoadBackedP256Key(cfg, tpm)
}
