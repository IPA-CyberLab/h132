package envelope

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    "envelope",
	Aliases: []string{"e"},
	Usage:   "Manage envelopes",
	Subcommands: []*cli.Command{
		sealCommand,
		unsealCommand,
		dumpCommand,
	},
}

func ReadFileCapped(fileName string, maxSize int64) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if fi.Size() > maxSize {
		return nil, fmt.Errorf("file size exceeds the maximum allowed size of %d bytes", maxSize)
	}

	bs := make([]byte, fi.Size())
	if _, err := io.ReadFull(f, bs); err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Verify that `f` reached EOF
	if _, err := f.Read(make([]byte, 1)); err != io.EOF {
		return nil, fmt.Errorf("file content was dynamically appended while reading")
	}

	return bs, nil
}
