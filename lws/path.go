package lws

import (
	"fmt"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
)

const ENVELOPE_FILEEXT = ".h132"

func GetLWSDir() string {
	d := os.Getenv("LWS_DIR")
	if d == "" {
		var err error
		d, err = os.Getwd()
		if err != nil {
			// Working directory is not available.
			// This is rare enough that we can panic here.
			zap.S().Fatalf("Failed to get working directory (probably it no longer exists): %v", err)
		}
	}

	return d
}

// CheckWriteAccess checks if the user has write access to the LWS directory.
func CheckWriteAccess() error {
	d := GetLWSDir()
	checkfn := path.Join(d, ".h132_write_access_check")

	f, err := os.OpenFile(checkfn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("access check failed - failed to create file %q: %w", checkfn, err)
	}
	f.Close()

	if err := os.Remove(checkfn); err != nil {
		return fmt.Errorf("access check failed - failed to remove file %q: %w", checkfn, err)
	}

	return nil
}

func GetLWSWireProtoPath() string {
	return path.Join(GetLWSDir(), "h132_letter_writing_set.binpb")
}

func GetEnvelopePath(plaintextPath string) string {
	fname := path.Base(plaintextPath) + ENVELOPE_FILEEXT
	return path.Join(GetLWSDir(), fname)
}

func GetPlaintextPath(envelopePath string) string {
	if strings.HasSuffix(envelopePath, ENVELOPE_FILEEXT) {
		return envelopePath[:len(envelopePath)-len(ENVELOPE_FILEEXT)]
	} else {
		return envelopePath + ".plaintext"
	}
}
