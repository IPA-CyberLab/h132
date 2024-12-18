package envelope

import (
	"errors"
	"fmt"
	"os"

	"github.com/IPA-CyberLab/h132/cmd/h132/common"
	"github.com/IPA-CyberLab/h132/cmd/h132/keys/access"
	"github.com/IPA-CyberLab/h132/lws"
	"github.com/urfave/cli/v2"
)

var editCommand = &cli.Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Usage:   "Edit h132 envelope file content",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "key",
			Usage:    "key name to be used for encryption and signing (use h132 keys list to see the list of available keys)",
			Required: true,
		},
		&cli.Int64Flag{
			Name:  "max-file-size",
			Usage: "maximum file size in bytes to be encrypted",
			Value: 1024 * 1024, // 1MB
		},
	},
	ArgsUsage: "FILE",
	Action: func(c *cli.Context) error {
		// s := zap.S()

		// FIXME[P2]: Run presubmit-like script to check if the latest change to the file is checked-in to git

		envelopePath := c.Args().First()
		if envelopePath == "" {
			return common.ErrInvalidInput{Msg: "file name is required"}
		}
		if c.Args().Len() > 1 {
			return common.ErrInvalidInput{Msg: "only one file name must be provided"}
		}

		maxFileSize := c.Int64("max-file-size")
		if maxFileSize == 0 {
			return common.ErrInvalidInput{Msg: "max-file-size must be greater than 0"}
		}

		l, err := lws.ReadLWS()
		if err != nil {
			return err
		}

		keyName := c.String("key")
		k := lws.GetKeyByName(l, keyName)
		if k == nil {
			return fmt.Errorf("specified key %q is not in the letter writing set %q", keyName, l.Name)
		}

		envelopeBs, err := ReadFileCapped(envelopePath, maxFileSize)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if err := lws.CheckWriteAccess(lws.GetLWSDir()); err != nil {
			return err
		}
		if err := lws.CheckWriteAccess(lws.GetPlaintextDir()); err != nil {
			return err
		}

		plaintextPath := lws.GetPlaintextPath(envelopePath)
		if err := lws.RunPreEditHook(l, envelopePath); err != nil {
			return err
		}

		ak, err := access.AccessKey(l.Name, k)
		if err != nil {
			return err
		}
		defer ak.Close()

		if err := lws.Edit(l, ak, envelopePath, envelopeBs, plaintextPath); err != nil {
			return err
		}

		return nil
	},
}
