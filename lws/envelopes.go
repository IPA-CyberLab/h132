package lws

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"

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
	s.Infof("Successfully sealed h132 envelope")

	if err := os.WriteFile(epath, buf.Bytes(), 0644); err != nil {
		return err
	}
	s.Infof("Successfully produced an envelope file %q", epath)

	return nil
}

func Unseal(ak envelope.AssymmetricKey, envelopeFilePath string, envelopeBs []byte) error {
	s := zap.S()

	plaintextPath := GetPlaintextPath(envelopeFilePath)

	r := bytes.NewReader(envelopeBs)
	plaintext, _, err := envelope.Unseal(r, ak)
	if err != nil {
		return err
	}
	s.Infof("Successfully unsealed h132 envelope %q", envelopeFilePath)

	if err := os.WriteFile(plaintextPath, plaintext, 0600); err != nil {
		return err
	}
	s.Infof("Successfully produced a plaintext file %q", plaintextPath)

	return nil
}

func RunEditor(filepath string) {
	s := zap.S()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		s.Infof("Failed to run $EDITOR %q on file %q: %v", editor, filepath, err)
	}
}

func Edit(l *pb.LetterWritingSet, ak envelope.AssymmetricKey, envelopePath string, envelopeBs []byte, plaintextPath string) error {
	s := zap.S()

	pubs := GetPublicKeys(l)
	if len(pubs) == 0 {
		return ErrNoKeys
	}

	var oldPlaintextBs []byte

	// Check if the plaintext file exists
	if _, err := os.Stat(plaintextPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to stat file %q: %w", plaintextPath, err)
		}
		// Proceed if the plaintext file does not exist.
	} else {
		var err error
		oldPlaintextBs, err = os.ReadFile(plaintextPath)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", plaintextPath, err)
		}
	}

	var plaintextBs []byte
	if len(envelopeBs) > 0 {
		var err error
		r := bytes.NewReader(envelopeBs)
		plaintextBs, _, err = envelope.Unseal(r, ak)
		if err != nil {
			return err
		}
	}

	if oldPlaintextBs != nil {
		if bytes.Equal(plaintextBs, oldPlaintextBs) {
			s.Infof("Found existing plaintext file %q with the same content as the envelope file %q. Proceeding.", plaintextPath, envelopePath)
		} else {
			return fmt.Errorf("the existing plaintext work file %q contains different content to the envelope file %q", plaintextPath, envelopePath)
		}
	}
	oldPlaintextBs = plaintextBs

	// Overwrite the plaintext file with the new content
	if err := os.WriteFile(plaintextPath, plaintextBs, 0600); err != nil {
		return err
	}

	RunEditor(plaintextPath)

	// Read the new content of the plaintext file
	plaintextBs, err := os.ReadFile(plaintextPath)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", plaintextPath, err)
	}

	if bytes.Equal(plaintextBs, oldPlaintextBs) {
		s.Infof("No change in the plaintext file %q. Skipping the sealing step.", plaintextPath)

		// Remove the plaintext file
		if err := os.Remove(plaintextPath); err != nil {
			return fmt.Errorf("failed to remove file %q: %w", plaintextPath, err)
		}
		s.Infof("Removed the plaintext file %q", plaintextPath)
		return nil
	}

	// Seal the new content
	var envBuf bytes.Buffer
	if err := envelope.Seal(&envBuf, plaintextBs, ak, pubs, rand.Reader); err != nil {
		return err
	}
	s.Infof("Successfully sealed h132 envelope")

	// Overwrite the envelope file with the new content
	if err := os.WriteFile(envelopePath, envBuf.Bytes(), 0644); err != nil {
		return err
	}
	s.Infof("Successfully produced an envelope file %q", envelopePath)

	// Remove the plaintext file
	if err := os.Remove(plaintextPath); err != nil {
		return fmt.Errorf("failed to remove file %q: %w", plaintextPath, err)
	}
	s.Infof("Removed the plaintext file %q", plaintextPath)

	if err := RunPostEditHook(l, envelopePath); err != nil {
		return err
	}

	return nil
}
