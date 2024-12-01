// The package must not be named "testdata" since go build tools ignore it.
package mytestdata

import (
	"crypto/ecdsa"
	"crypto/x509"
	"embed"
	"encoding/pem"
	"log"
)

//go:embed priv1.pem priv2.pem priv3.pem
var embedFs embed.FS

var (
	Priv1 *ecdsa.PrivateKey
	Priv2 *ecdsa.PrivateKey
	Priv3 *ecdsa.PrivateKey
)

func readPrivateKey(filename string) *ecdsa.PrivateKey {
	data, err := embedFs.ReadFile(filename)
	if err != nil {
		log.Panicf("readPrivateKey: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		log.Panicf("pem.Decode: %v", err)
	}

	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		log.Panicf("x509.ParseECPrivateKey: %v", err)
	}

	return priv
}

func init() {
	Priv1 = readPrivateKey("priv1.pem")
	Priv2 = readPrivateKey("priv2.pem")
	Priv3 = readPrivateKey("priv3.pem")
}
