package access

import (
	"fmt"

	"github.com/IPA-CyberLab/h132/cmd/h132/keys/access/promptmneu"
	"github.com/IPA-CyberLab/h132/cmd/h132/promptcode"
	"github.com/IPA-CyberLab/h132/envelope"
	"github.com/IPA-CyberLab/h132/keys/emergency"
	"github.com/IPA-CyberLab/h132/keys/webauthnwrappedtpm"
	"github.com/IPA-CyberLab/h132/pb"
	h132_tpm2 "github.com/IPA-CyberLab/h132/tpm2"
	"github.com/IPA-CyberLab/h132/webauthn"
)

// AccessKey retrieves the envelope decrypting/signing key by interacting with the user.
func AccessKey(lwsName string, k *pb.KeyImpl) (envelope.AssymmetricKey, error) {
	pub := pb.ProtoToPub(k.PublicKey)

	switch ki := k.Impl.(type) {
	case *pb.KeyImpl_EmergencyBackupKey:
		hint := ki.EmergencyBackupKey.Hint

		mneu, err := promptmneu.Prompt(k.Name, hint)
		if err != nil {
			return nil, err
		}

		ek, err := emergency.FromMneumonic(hint, pub, mneu)
		if err != nil {
			return nil, err
		}

		return ek, nil

	case *pb.KeyImpl_WebauthnWrappedTpm:
		tpm, err := h132_tpm2.GetTPM()
		if err != nil {
			return nil, fmt.Errorf("failed to get TPM: %w", err)
		}
		defer func() {
			if tpm != nil {
				tpm.Close()
			}
		}()

		wwt := ki.WebauthnWrappedTpm

		wk, err := webauthnwrappedtpm.GetProvisioned(tpm, lwsName, k.Name, wwt)
		if err != nil {
			return nil, err
		}

		prfsess, err := wk.StartGetPRFSecretSession()
		if err != nil {
			return nil, err
		}

		var b64resp []byte
		if err := promptcode.Prompt(prfsess.GetPRFSecretUrlStr, func(code string) error {
			bs := []byte(code)
			if _, err := webauthn.B64ToBytes(bs); err != nil {
				return err
			}

			b64resp = bs
			return nil
		}); err != nil {
			return nil, err
		}
		prf, err := prfsess.Complete(b64resp)
		if err != nil {
			return nil, err
		}

		bk, err := wk.Unwrap(tpm, prf)
		if err != nil {
			return nil, err
		}
		// `bk` has taken over ownership of tpm. Don't close it here.
		tpm = nil

		return bk, nil

	default:
		return nil, fmt.Errorf("unknown key type: %T", ki)
	}
}
