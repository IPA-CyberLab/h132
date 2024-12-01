package keys

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/IPA-CyberLab/h132/cmd/h132/keys/promptnewemergency"
	"github.com/IPA-CyberLab/h132/cmd/h132/promptcode"
	"github.com/IPA-CyberLab/h132/keys/emergency"
	"github.com/IPA-CyberLab/h132/keys/webauthnwrappedtpm"
	"github.com/IPA-CyberLab/h132/lws"
	"github.com/IPA-CyberLab/h132/pb"
	h132_tpm2 "github.com/IPA-CyberLab/h132/tpm2"
	"github.com/IPA-CyberLab/h132/webauthn"
)

func addEmergencyKey(c *cli.Context, name, hint string) error {
	s := zap.S()

	l, err := lws.ReadLWS()
	if err != nil {
		return err
	}
	if lws.GetKeyByName(l, name) != nil {
		return fmt.Errorf("key %q already exists in the letter writing set %q", name, l.Name)
	}

	newFunc := func() *emergency.EmergencyBackupKey { return emergency.New(hint, rand.Reader) }
	ek, err := promptnewemergency.Prompt(newFunc)
	if err != nil {
		return err
	}

	keyproto, err := pb.KeyToProto(name, ek)
	if err != nil {
		return err
	}
	s.Infof("Successfully generated key: %s", pb.KeyImplSummary(keyproto))

	l.Keys = append(l.Keys, keyproto)
	if err := lws.UpdateLWS(l, 0); err != nil {
		return fmt.Errorf("failed to update lws: %w", err)
	}

	s.Infof("Successfully added the key to the letter writing set.")
	return nil
}

func addWebauthnWrappedTPMKey(c *cli.Context, name, reflectorUrlStr string, tpm transport.TPM, tpmKeyHandle uint32) error {
	s := zap.S()

	l, err := lws.ReadLWS()
	if err != nil {
		return err
	}
	if lws.GetKeyByName(l, name) != nil {
		return fmt.Errorf("key %q already exists in the letter writing set %q", name, l.Name)
	}

	wk, err := webauthnwrappedtpm.NewUnprovisioned(tpm, reflectorUrlStr, tpm2.TPMHandle(tpmKeyHandle))
	if err != nil {
		return err
	}

	sess, err := wk.StartRegistration(l.Name, name)
	if err != nil {
		return err
	}

	var zb64resp []byte
	if err := promptcode.Prompt(sess.RegistrationUrlStr, func(code string) error {
		bs := []byte(code)
		if _, err := webauthn.Zb64ToBytes(bs); err != nil {
			return err
		}

		zb64resp = bs
		return nil
	}); err != nil {
		return err
	}

	if err := wk.CompleteRegistration(sess, zb64resp); err != nil {
		return err
	}

	s.Infof("Successfully registered webauthn credential. Now onto acquiring PRF secret.")

	prfsess, err := wk.StartGetPRFSecretSession()
	if err != nil {
		return err
	}

	if err := promptcode.Prompt(prfsess.GetPRFSecretUrlStr, func(code string) error {
		bs := []byte(code)
		if _, err := webauthn.B64ToBytes(bs); err != nil {
			return err
		}

		zb64resp = bs
		return nil
	}); err != nil {
		return err
	}
	prf, err := prfsess.Complete(zb64resp)
	if err != nil {
		return err
	}

	s.Info("Successfully acquired prf. Now provisioning the key on the TPM.")
	if err := wk.Provision(tpm, prf); err != nil {
		return err
	}

	keyproto, err := pb.KeyToProto(name, wk)
	if err != nil {
		return err
	}
	s.Infof("Successfully generated key: %s", pb.KeyImplSummary(keyproto))

	l.Keys = append(l.Keys, keyproto)
	if err := lws.UpdateLWS(l, 0); err != nil {
		return fmt.Errorf("failed to update lws: %w", err)
	}

	s.Infof("Successfully added the key to the letter writing set.")
	return nil
}

var addCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "Add a key",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "type",
			Usage:    "key type (emergency, webauthn_wrapped_tpm)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "name of the key to be added",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "hint",
			Usage: "Hint about the physical location of the key mneumonic written down.",
		},
		&cli.StringFlag{
			Name:  "reflectorUrl",
			Usage: "URL of the reflector server",
			Value: "https://h132-fixme.example/",
		},
		&cli.StringFlag{
			Name:  "tpmKeyHandle",
			Usage: "TPM key handle in hex (e.g. 81008001)",
		},
	},
	Action: func(c *cli.Context) error {
		name := c.String("name")
		if name == "" {
			return fmt.Errorf("name is required")
		}
		if strings.Contains(name, "@") {
			return fmt.Errorf("name cannot contain '@'")
		}

		typ := c.String("type")
		switch typ {
		case "emergency":
			hint := c.String("hint")
			if hint == "" {
				return fmt.Errorf("hint is required")
			}

			return addEmergencyKey(c, name, hint)
		case "webauthn_wrapped_tpm":
			reflectorUrlStr := c.String("reflectorUrl")
			if reflectorUrlStr == "" {
				return fmt.Errorf("reflectorUrl cannot be empty")
			}

			tkh := c.String("tpmKeyHandle")
			if tkh == "" {
				return fmt.Errorf("tpmKeyHandle is required")
			}
			tkh = strings.TrimPrefix(tkh, "0x")
			tkhu, err := strconv.ParseUint(tkh, 16, 32)
			if err != nil {
				return fmt.Errorf("failed to parse tpmKeyHandle: %w", err)
			}

			tpm, err := h132_tpm2.GetTPM()
			if err != nil {
				return fmt.Errorf("failed to get TPM: %w", err)
			}
			defer tpm.Close()

			return addWebauthnWrappedTPMKey(c, name, reflectorUrlStr, tpm, uint32(tkhu))
		default:
			return fmt.Errorf("unknown key type: %s", typ)
		}
	},
}
