package envelope

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	testdata "github.com/IPA-CyberLab/h132/mytestdata"
)

func TestEnvelope_Basic(t *testing.T) {
	plaintext := []byte("himichu")

	akey := NewLocalPrivateKey(testdata.Priv1)
	if !akey.Public().Equal(&testdata.Priv1.PublicKey) {
		t.Error("akey.Public() != testdata.Priv1.PublicKey")
	}

	var buf bytes.Buffer
	err := Seal(
		&buf, plaintext, akey,
		[]*ecdsa.PublicKey{&testdata.Priv1.PublicKey}, rand.Reader)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}

	plaintext2, recs, err := Unseal(&buf, akey)
	if err != nil {
		t.Fatalf("Unseal: %v", err)
	}
	if !bytes.Equal(plaintext, plaintext2) {
		t.Fatalf("plaintext != plaintext2")
	}
	if len(recs) != 1 {
		t.Fatalf("len(recs) = %d", len(recs))
	}
	if !recs[0].Equal(&testdata.Priv1.PublicKey) {
		t.Fatalf("recs[0] != testdata.Priv1.PublicKey")
	}
}
