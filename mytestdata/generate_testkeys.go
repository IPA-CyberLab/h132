//go:build ignore

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func genECPrivateKeyIfNotExists(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Printf("file %s already exists", filename)
		return nil
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("ecdsa.GenerateKey: %w", err)
	}

	kbs, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return fmt.Errorf("x509.MarshalECPrivateKey: %w", err)
	}

	bs := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: kbs,
	})
	if err := os.WriteFile(filename, bs, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}
	log.Printf("file %s created", filename)
	return nil
}

func run() error {
	if err := genECPrivateKeyIfNotExists("priv1.pem"); err != nil {
		return err
	}
	if err := genECPrivateKeyIfNotExists("priv2.pem"); err != nil {
		return err
	}
	if err := genECPrivateKeyIfNotExists("priv3.pem"); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Printf("%v", err)
	}
}
