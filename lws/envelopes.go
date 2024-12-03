package lws

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"os"

	"github.com/IPA-CyberLab/h132/envelope"
	"github.com/IPA-CyberLab/h132/pb"
	"go.uber.org/zap"
)

var (
	ErrNoKeys = errors.New("no keys in the letter writing set")
)

func GetPublicKeys(l *pb.LetterWritingSet) []*ecdsa.PublicKey {
	pubs := make([]*ecdsa.PublicKey, 0, len(l.Keys))
	for _, k := range l.Keys {
		pubs = append(pubs, pb.ProtoToPub(k.PublicKey))
	}
	return pubs
}

func Seal(l *pb.LetterWritingSet, ak envelope.AssymmetricKey, fileName string, contents []byte) error {
	s := zap.S()

	epath := GetEnvelopePath(fileName)

	pubs := GetPublicKeys(l)
	if len(pubs) == 0 {
		return ErrNoKeys
	}

	var buf bytes.Buffer
	if err := envelope.Seal(&buf, contents, ak, pubs, rand.Reader); err != nil {
		return err
	}
	s.Infof("Successfully sealed h132 envelope %q", epath)

	if err := os.WriteFile(epath, buf.Bytes(), 0644); err != nil {
		return err
	}
	s.Infof("Successfully produced an envelope file %q", epath)

	return nil
}

func Unseal(ak envelope.AssymmetricKey, envelopeFileName string, envelopeBs []byte) error {
	s := zap.S()

	epath := GetPlaintextPath(envelopeFileName)

	r := bytes.NewReader(envelopeBs)
	plaintext, _, err := envelope.Unseal(r, ak)
	if err != nil {
		return err
	}
	s.Infof("Successfully unsealed h132 envelope %q", envelopeFileName)

	if err := os.WriteFile(epath, plaintext, 0600); err != nil {
		return err
	}
	s.Infof("Successfully produced a plaintext file %q", epath)

	return nil
}
