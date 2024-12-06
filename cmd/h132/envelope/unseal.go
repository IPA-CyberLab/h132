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

var unsealCommand = &cli.Command{
	Name:    "unseal",
	Aliases: []string{"u"},
	Usage:   "Decrypt a h132 envelope file and restore the original file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "key",
			Usage:    "key name to be used for decryption (use h132 keys list to see the list of available keys)",
			Required: true,
		},
	},
	ArgsUsage: "FILE",
	Action: func(c *cli.Context) error {
		// s := zap.S()

		envelopeFileName := c.Args().First()
		if envelopeFileName == "" {
			return common.ErrInvalidInput{Msg: "file name is required"}
		}
		if c.Args().Len() > 1 {
			return common.ErrInvalidInput{Msg: "only one file name must be provided"}
		}

		bs, err := os.ReadFile(envelopeFileName)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", envelopeFileName, err)
		}

		l, err := lws.ReadLWS()
		if err != nil {
			return err
		}
		// Note: we output the decrypted file to the LWS dir, so we need write access
		if err := lws.CheckWriteAccess(lws.GetLWSDir()); err != nil {
			return err
		}

		keyName := c.String("key")
		k := lws.GetKeyByName(l, keyName)
		if k == nil {
			return fmt.Errorf("specified key %q is not in the letter writing set %q", keyName, l.Name)
		}

		outfp := lws.GetPlaintextPath(envelopeFileName)
		fi, err := os.Stat(outfp)
		if errors.Is(err, os.ErrNotExist) {
			// File does not exist, so we can proceed
		} else if err != nil {
			return fmt.Errorf("failed to stat file %q: %w", outfp, err)
		} else if fi.IsDir() {
			return fmt.Errorf("file %q is a directory", outfp)
		} else {
			return fmt.Errorf("file %q: %w", outfp, os.ErrExist)
		}

		ak, err := access.AccessKey(l.Name, k)
		if err != nil {
			return err
		}
		defer ak.Close()

		if err := lws.Unseal(ak, envelopeFileName, bs); err != nil {
			return fmt.Errorf("failed to unseal file %q: %w", envelopeFileName, err)
		}
		return nil
	},
}
