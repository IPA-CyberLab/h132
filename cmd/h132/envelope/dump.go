package envelope

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/IPA-CyberLab/h132/cmd/h132/common"
	"github.com/IPA-CyberLab/h132/envelope"
	"github.com/IPA-CyberLab/h132/lws"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var dumpCommand = &cli.Command{
	Name:    "dump",
	Aliases: []string{"d"},
	Usage:   "Dump the content of a h132 envelope file for debugging purposes",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Dump the hexdump of ciphertext.",
		},
	},
	ArgsUsage: "FILE.h132",
	Action: func(c *cli.Context) error {
		s := zap.S()

		fileName := c.Args().First()
		if fileName == "" {
			return common.ErrInvalidInput{Msg: "file name is required"}
		}
		if c.Args().Len() > 1 {
			return common.ErrInvalidInput{Msg: "only one file name must be provided"}
		}

		f, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()

		annotator := func(*ecdsa.PublicKey) string { return "" }

		l, err := lws.ReadLWS()
		if err != nil {
			s.Error("Failed to read the letter writing set. Continuing without annotating public keys.")
		} else {
			annotator = func(pub *ecdsa.PublicKey) string {
				key := lws.GetKeyByPublicKey(l, pub)
				if key == nil {
					return "[key not in the letter writing set]"
				}
				return key.Name
			}
		}

		dumpstr, err := envelope.Dump(f, c.Bool("verbose"), annotator)
		fmt.Fprintf(os.Stdout, "%s", dumpstr)
		if err != nil {
			return fmt.Errorf("Error: %w", err)
		}

		return nil
	},
}
