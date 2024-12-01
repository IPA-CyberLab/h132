package lws_test

import (
	"os"
	"testing"

	"github.com/IPA-CyberLab/h132/lws"
)

func TestCheckWriteAccess(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "lws_test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set the LWS_DIR environment variable to the temporary directory
	os.Setenv("LWS_DIR", tempDir)
	defer os.Unsetenv("LWS_DIR")

	// Test case where write access is available
	err = lws.CheckWriteAccess()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test case where write access is not available
	// Remove write permissions from the temporary directory
	if err := os.Chmod(tempDir, 0500); err != nil {
		t.Fatalf("failed to change directory permissions: %v", err)
	}

	err = lws.CheckWriteAccess()
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// Restore write permissions to the temporary directory
	if err := os.Chmod(tempDir, 0700); err != nil {
		t.Fatalf("failed to restore directory permissions: %v", err)
	}
}

func TestGetEnvelopePath(t *testing.T) {
	os.Setenv("LWS_DIR", "/foo/bar/lws")
	defer os.Unsetenv("LWS_DIR")

	tcs := []struct {
		plaintextPath  string
		expectedResult string
	}{
		{
			plaintextPath:  "/tmp/letter",
			expectedResult: "/foo/bar/lws/letter.h132",
		},
		{
			plaintextPath:  "/tmp/another_letter.txt",
			expectedResult: "/foo/bar/lws/another_letter.txt.h132",
		},
	}
	for _, tc := range tcs {
		result := lws.GetEnvelopePath(tc.plaintextPath)
		if result != tc.expectedResult {
			t.Errorf("expected %q, got %q", tc.expectedResult, result)
		}
	}
}

func TestGetPlaintextPath(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		envelopePath   string
		expectedResult string
	}{
		{
			envelopePath:   "/tmp/letter.h132",
			expectedResult: "/tmp/letter",
		},
		{
			envelopePath:   "/tmp/letter",
			expectedResult: "/tmp/letter.plaintext",
		},
	}
	for _, tc := range tcs {
		result := lws.GetPlaintextPath(tc.envelopePath)
		if result != tc.expectedResult {
			t.Errorf("expected %q, got %q", tc.expectedResult, result)
		}
	}
}
