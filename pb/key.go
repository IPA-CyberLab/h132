package pb

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"log"
	"math/big"

	"go.uber.org/zap"
)

func PubToProto(pub *ecdsa.PublicKey) *P256PublicKey {
	if pub.Curve.Params() != elliptic.P256().Params() {
		log.Panicf("PubToProto: unsupported curve %v", pub.Curve)
	}

	return &P256PublicKey{
		X: pub.X.Bytes(),
		Y: pub.Y.Bytes(),
	}
}

func ProtoToPub(proto *P256PublicKey) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(proto.X),
		Y:     new(big.Int).SetBytes(proto.Y),
	}
}

func PubKeyPinString(pub *ecdsa.PublicKey) string {
	pkix, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		zap.S().Panicf("Failed to marshal public key: %v", err)
		return ""
	}

	hash := sha256.Sum256(pkix)
	str := base64.StdEncoding.EncodeToString(hash[:])
	return str
}

func PrivToProto(priv ecdsa.PrivateKey) *P256PrivateKey {
	if priv.Curve.Params() != elliptic.P256().Params() {
		log.Panicf("PrivToProto: unsupported curve %v", priv.Curve)
	}

	return &P256PrivateKey{
		D: priv.D.Bytes(),
	}
}

func ProtoToPriv(proto *P256PrivateKey) ecdsa.PrivateKey {
	priv := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     nil,
			Y:     nil,
		},
		D: new(big.Int).SetBytes(proto.D),
	}

	// Note: `ScalarBaseMult` is deprecated without an alternative.
	priv.PublicKey.X, priv.PublicKey.Y =
		priv.PublicKey.Curve.ScalarBaseMult(priv.D.Bytes())

	return priv
}
