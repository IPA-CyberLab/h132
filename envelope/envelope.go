package envelope

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"

	"go.uber.org/multierr"
	"golang.org/x/crypto/hkdf"
	"google.golang.org/protobuf/proto"

	"github.com/IPA-CyberLab/h132/pb"
)

var (
	EncryptSymmetricKeySalt = []byte("encrypt_symmetric_key")

	ErrInvalidEnvelopeSignature     = errors.New("Envelope signature verification failed.")
	ErrInvalidRecipientKeySignature = errors.New("Recipient key signature verification failed.")
	ErrNotRecipient                 = errors.New("The private key is not a recipient of the envelope.")
)

func deriveSymmetricKeyEncryptionKey(sharedSecret []byte) []byte {
	hkdfr := hkdf.New(sha256.New, sharedSecret, EncryptSymmetricKeySalt, []byte("message_symmetric_key"))
	key := make([]byte, 32)
	if _, err := hkdfr.Read(key); err != nil {
		log.Fatalf("hkdfr.Read: %v", err)
	}
	return key
}

func EncryptSymmetricKey(symmetricKey []byte, akey AssymmetricKey, recipient *ecdsa.PublicKey, randr io.Reader) (*pb.EncryptedSymmetricKey, error) {
	// Generate ephemeral key pair
	ephPriv, err := ecdsa.GenerateKey(elliptic.P256(), randr)
	if err != nil {
		return nil, fmt.Errorf("ecdsa.GenerateKey: %w", err)
	}

	// Sign the ephemeral public key so the recipient can verify it
	hash := sha256.New()
	hash.Write(ephPriv.PublicKey.X.Bytes())
	hash.Write(ephPriv.PublicKey.Y.Bytes())
	ephPubDigest := hash.Sum(nil)

	publicKeySign, err := akey.Sign(ephPubDigest)
	if err != nil {
		return nil, fmt.Errorf("priv.Sign(ephPubDigest): %w", err)
	}

	// Derive shared secret using ECDH
	privDH, err := ephPriv.ECDH()
	if err != nil {
		return nil, fmt.Errorf("priv.ECDH: %w", err)
	}

	pubDH, err := recipient.ECDH()
	if err != nil {
		return nil, fmt.Errorf("recipient.ECDH: %w", err)
	}

	sharedSecret, err := privDH.ECDH(pubDH)
	if err != nil {
		return nil, fmt.Errorf("privDH.ECDH(pubDH): %w", err)
	}

	// log.Printf("  encSymKey - sharedSecret: %x", sharedSecret)

	// Use HKDF to derive a symmetric key for encryptng the `key`
	keyEncKey := deriveSymmetricKeyEncryptionKey(sharedSecret)

	// Encrypt the `key` with the derived symmetric key `keyEncKey`
	aesc, err := aes.NewCipher(keyEncKey)
	if err != nil {
		log.Fatalf("aes.NewCipher(encKey): %v", err)
	}
	aead, err := cipher.NewGCM(aesc)
	if err != nil {
		log.Fatalf("cipher.NewGCM(aesc): %v", err)
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := randr.Read(nonce); err != nil {
		return nil, fmt.Errorf("rand.Read(nonce): %w", err)
	}

	encKey := aead.Seal(nil, nonce, symmetricKey, nil)

	return &pb.EncryptedSymmetricKey{
		RecipientPublicKey:           pb.PubToProto(recipient),
		SenderEphemeralPublicKey:     pb.PubToProto(&ephPriv.PublicKey),
		SenderEphemeralPublicKeySign: publicKeySign,
		EncryptedSymmetricKey:        encKey,
		Nonce:                        nonce,
	}, nil
}

// Verify the sender ephemeral public key signature
func verifySenderEphemeralPublicKeySign(encryptedKey *pb.EncryptedSymmetricKey, senderPub *ecdsa.PublicKey) bool {
	hash := sha256.New()
	hash.Write(encryptedKey.SenderEphemeralPublicKey.X)
	hash.Write(encryptedKey.SenderEphemeralPublicKey.Y)
	ephPubDigest := hash.Sum(nil)

	return ecdsa.VerifyASN1(senderPub, ephPubDigest, encryptedKey.SenderEphemeralPublicKeySign)
}

// Decrypt the encrypted symmetric key `encryptedKey` using the private key `akey`
func decryptSymmetricKey(encryptedKey *pb.EncryptedSymmetricKey, akey AssymmetricKey, senderPub *ecdsa.PublicKey) ([]byte, error) {
	if !verifySenderEphemeralPublicKeySign(encryptedKey, senderPub) {
		return nil, ErrInvalidRecipientKeySignature
	}

	// Derive shared secret using ECDH
	ephPub := pb.ProtoToPub(encryptedKey.SenderEphemeralPublicKey)
	sharedSecret, err := akey.ECDH(ephPub)
	if err != nil {
		return nil, fmt.Errorf("privDH.ECDH(pubDH): %w", err)
	}

	// log.Printf("  decSymKey - sharedSecret: %x", sharedSecret)

	// Use HKDF to derive the symmetric key used to encrypt the `key`
	keyEncKey := deriveSymmetricKeyEncryptionKey(sharedSecret)

	// Decrypt the `key` with the derived symmetric key `keyEncKey`
	aesc, err := aes.NewCipher(keyEncKey)
	if err != nil {
		log.Fatalf("aes.NewCipher(encKey): %v", err)
	}
	aead, err := cipher.NewGCM(aesc)
	if err != nil {
		log.Fatalf("cipher.NewGCM(aesc): %v", err)
	}

	key, err := aead.Open(nil, encryptedKey.Nonce, encryptedKey.EncryptedSymmetricKey, nil)
	if err != nil {
		return nil, fmt.Errorf("aead.Open: %w", err)
	}

	return key, nil
}

func Seal(w io.Writer, plaintext []byte, akey AssymmetricKey, recipients []*ecdsa.PublicKey, randr io.Reader) error {
	// Generate random symmetric key used to encrypt the `plaintext` message
	symmetricKey := make([]byte, 32)
	if _, err := randr.Read(symmetricKey); err != nil {
		return fmt.Errorf("rand.Read(symmetricKey): %w", err)
	}

	// log.Printf("  seal - symmetricKey: %x", symmetricKey)

	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return fmt.Errorf("aes.NewCipher: %w", err)
	}

	// Encrypt the `plaintext` message with the symmetric key
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("aes.NewGCM: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = randr.Read(nonce); err != nil {
		return fmt.Errorf("rand.Read(nonce): %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	// Encrypt the symmetric key with the public key of each recipient
	recipientKeys := make([]*pb.EncryptedSymmetricKey, len(recipients))
	for i, recipient := range recipients {
		rk, err := EncryptSymmetricKey(symmetricKey, akey, recipient, randr)
		if err != nil {
			return fmt.Errorf("EncryptSymmetricKey: %w", err)
		}
		recipientKeys[i] = rk
	}

	lr := &pb.Letter{
		Ciphertext:      ciphertext,
		Nonce:           nonce,
		RecipientKeys:   recipientKeys,
		SenderPublicKey: pb.PubToProto(akey.Public()),
	}
	lrbs, err := proto.Marshal(lr)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	digest := sha256.Sum256(lrbs)
	signature, err := akey.Sign(digest[:])
	if err != nil {
		return fmt.Errorf("priv.Sign(envelope_digest): %w", err)
	}

	ev := &pb.Envelope{
		LetterProto: lrbs,
		Signature:   signature,
	}
	evbs, err := proto.Marshal(ev)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}
	if _, err = w.Write(evbs); err != nil {
		return fmt.Errorf("w.Write(evbs): %w", err)
	}

	return nil
}

func read(r io.Reader) (*pb.Envelope, *pb.Letter, error) {
	evbs, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	var ev pb.Envelope
	if err := proto.Unmarshal(evbs, &ev); err != nil {
		return nil, nil, fmt.Errorf("proto.UnmarshalReader(pb.Envelope): %w", err)
	}
	var lr pb.Letter
	if err := proto.Unmarshal(ev.LetterProto, &lr); err != nil {
		return nil, nil, fmt.Errorf("proto.UnmarshalReader(pb.Letter): %w", err)
	}

	return &ev, &lr, nil
}

type PublicKeyAnnotatorFunc func(*ecdsa.PublicKey) string

func Dump(r io.Reader, verbose bool, pubaf PublicKeyAnnotatorFunc) (string, error) {
	var out bytes.Buffer

	ev, lr, err := read(r)
	if err != nil {
		return "!!! Failed to read envelope binpb", err
	}

	out.WriteString(("Envelope:\n"))
	fmt.Fprintf(&out, "  Signature: %x", ev.Signature)

	digest := sha256.Sum256(ev.LetterProto)

	senderPub := pb.ProtoToPub(lr.SenderPublicKey)
	if ok := ecdsa.VerifyASN1(senderPub, digest[:], ev.Signature); ok {
		out.WriteString(" (valid)\n")
	} else {
		multierr.AppendInto(&err, ErrInvalidEnvelopeSignature)
		out.WriteString(" (invalid)\n")
	}

	out.WriteString("  Letter:\n")
	fmt.Fprintf(&out, "    len(Ciphertext): %d\n", len(lr.Ciphertext))
	if verbose {
		fmt.Fprintf(&out, "    Ciphertext:\n%s", hex.Dump(lr.Ciphertext))
	}
	fmt.Fprintf(&out, "    Nonce: %s\n", hex.EncodeToString(lr.Nonce))
	for i, rk := range lr.RecipientKeys {
		out.WriteString(fmt.Sprintf("    RecipientKey[%d]:\n", i))
		rpub := pb.ProtoToPub(rk.RecipientPublicKey)
		fmt.Fprintf(&out, "      RecipientPublicKey: %s", pb.PubKeyPinString(rpub))
		if rpubA := pubaf(rpub); rpubA != "" {
			fmt.Fprintf(&out, " (%s)", rpubA)
		}
		out.WriteString("\n")
		fmt.Fprintf(&out, "      SenderEphemeralPublicKey: %s\n", pb.PubKeyPinString(pb.ProtoToPub(rk.SenderEphemeralPublicKey)))
		fmt.Fprintf(&out, "      SenderEphemeralPublicKeySign: %x", rk.SenderEphemeralPublicKeySign)
		if ok := verifySenderEphemeralPublicKeySign(rk, senderPub); ok {
			out.WriteString(" (valid)\n")
		} else {
			multierr.AppendInto(&err, ErrInvalidRecipientKeySignature)
			out.WriteString(" (invalid)\n")
		}
		fmt.Fprintf(&out, "      EncryptedSymmetricKey: %s\n", hex.EncodeToString(rk.EncryptedSymmetricKey))
		fmt.Fprintf(&out, "      Nonce: %s\n", hex.EncodeToString(rk.Nonce))
	}
	spub := pb.ProtoToPub(lr.SenderPublicKey)
	fmt.Fprintf(&out, "    SenderPublicKey: %s", pb.PubKeyPinString(spub))
	if spubA := pubaf(spub); spubA != "" {
		fmt.Fprintf(&out, " (%s)", spubA)
	}
	out.WriteString("\n")

	return out.String(), err
}

func Unseal(r io.Reader, akey AssymmetricKey) ([]byte, []*ecdsa.PublicKey, error) {
	ev, lr, err := read(r)
	if err != nil {
		return nil, nil, err
	}

	// Verify the letter signature
	digest := sha256.Sum256(ev.LetterProto)

	senderPub := pb.ProtoToPub(lr.SenderPublicKey)
	if ok := ecdsa.VerifyASN1(senderPub, digest[:], ev.Signature); !ok {
		return nil, nil, ErrInvalidEnvelopeSignature
	}

	// Construct the array of recipient public keys, and find the recipient
	// key that matches the private key.
	akeyPub := akey.Public()
	recipients := make([]*ecdsa.PublicKey, len(lr.RecipientKeys))
	var myKey *pb.EncryptedSymmetricKey
	for i, rk := range lr.RecipientKeys {
		pub := pb.ProtoToPub(rk.RecipientPublicKey)
		recipients[i] = pub

		if pub.Equal(akeyPub) {
			myKey = rk
		}
	}
	if myKey == nil {
		return nil, recipients, ErrNotRecipient
	}

	symmetricKey, err := decryptSymmetricKey(myKey, akey, senderPub)
	if err != nil {
		return nil, recipients, err
	}

	// log.Printf("unseal - symmetricKey: %x", symmetricKey)

	// Decrypt the ciphertext using the `symmetricKey`
	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, recipients, fmt.Errorf("aes.NewCipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, recipients, fmt.Errorf("aes.NewGCM: %w", err)
	}

	plaintext, err := aead.Open(nil, lr.Nonce, lr.Ciphertext, nil)
	if err != nil {
		return nil, recipients, fmt.Errorf("aead.Open: %w", err)
	}

	return plaintext, recipients, nil
}
