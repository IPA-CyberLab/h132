package envelope

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io"
)

type AssymmetricKey interface {
	Public() *ecdsa.PublicKey
	Sign(s256digest []byte) ([]byte, error)
	ECDH(pub *ecdsa.PublicKey) ([]byte, error)
	io.Closer
}

type LocalPrivateKey struct {
	randr io.Reader
	Priv  *ecdsa.PrivateKey
}

var _ AssymmetricKey = &LocalPrivateKey{}

func NewLocalPrivateKey(priv *ecdsa.PrivateKey) *LocalPrivateKey {
	return &LocalPrivateKey{randr: rand.Reader, Priv: priv}
}

func (k *LocalPrivateKey) Public() *ecdsa.PublicKey {
	return &k.Priv.PublicKey
}

func (k *LocalPrivateKey) Sign(s256digest []byte) ([]byte, error) {
	if len(s256digest) != 32 {
		return nil, fmt.Errorf("LocalPrivateKey.Sign(): invalid digest length")
	}

	return k.Priv.Sign(k.randr, s256digest, nil)
}

func (k *LocalPrivateKey) ECDH(pub *ecdsa.PublicKey) ([]byte, error) {
	privDH, err := k.Priv.ECDH()
	if err != nil {
		return nil, fmt.Errorf("(ecdsa.PrivateKey).ECDH(): %w", err)
	}

	pubDH, err := pub.ECDH()
	if err != nil {
		return nil, fmt.Errorf("recipient.ECDH: %w", err)
	}

	return privDH.ECDH(pubDH)
}

func (k *LocalPrivateKey) Close() error {
	return nil
}
