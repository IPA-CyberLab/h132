package emergency

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"io"
	"math/big"

	"github.com/IPA-CyberLab/h132/envelope"
	"github.com/IPA-CyberLab/h132/pb"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"go.uber.org/zap"
)

type EmergencyBackupKey struct {
	*envelope.LocalPrivateKey

	// Hint about the physical location of the key mneumonic written down.
	hint string
}

func New(hint string, randr io.Reader) *EmergencyBackupKey {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), randr)
	if err != nil {
		zap.S().Panicf("Failed to generate ECDSA key: %v", err)
	}

	return &EmergencyBackupKey{
		// FIXME[?]: `randr` is not passed to LocalPrivateKey
		LocalPrivateKey: envelope.NewLocalPrivateKey(priv),
		hint:            hint,
	}
}

func (k *EmergencyBackupKey) Mneumonic() string {
	dbytes := k.LocalPrivateKey.Priv.D.Bytes()

	bip39.SetWordList(wordlists.Japanese)
	m, err := bip39.NewMnemonic(dbytes)
	if err != nil {
		zap.S().Panicf("Failed to generate mnemonic: %v", err)
	}

	return m
}

func FromMneumonic(hint string, pub *ecdsa.PublicKey, mneumonic string) (*EmergencyBackupKey, error) {
	dbytes, err := bip39.EntropyFromMnemonic(mneumonic)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mneumonic: %w", err)
	}

	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     nil,
			Y:     nil,
		},
		D: new(big.Int).SetBytes(dbytes),
	}
	// Note: `ScalarBaseMult` is deprecated without an alternative.
	priv.PublicKey.X, priv.PublicKey.Y =
		priv.PublicKey.Curve.ScalarBaseMult(priv.D.Bytes())

	if pub.X.Cmp(priv.PublicKey.X) != 0 || pub.Y.Cmp(priv.PublicKey.Y) != 0 {
		return nil, fmt.Errorf("private key recovered from mneumonic does not match the recorded public key")
	}

	return &EmergencyBackupKey{
		LocalPrivateKey: envelope.NewLocalPrivateKey(priv),
		hint:            hint,
	}, nil
}

func (k *EmergencyBackupKey) Pub() (*ecdsa.PublicKey, error) {
	return &k.Priv.PublicKey, nil
}

func (k *EmergencyBackupKey) SetImplProto(p *pb.KeyImpl) error {
	p.Impl = &pb.KeyImpl_EmergencyBackupKey{
		EmergencyBackupKey: &pb.EmergencyBackupKey{
			Hint: k.hint,
		},
	}

	return nil
}
