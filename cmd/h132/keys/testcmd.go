package keys

import (
	"crypto/sha256"
	"fmt"

	"github.com/IPA-CyberLab/h132/cmd/h132/keys/access"
	"github.com/IPA-CyberLab/h132/lws"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var testCommand = &cli.Command{
	Name:    "test",
	Aliases: []string{"t"},
	Usage:   "Test the key access procedure and confirm public key",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "name of the key to be tested",
		},
	},
	Action: func(c *cli.Context) error {
		s := zap.S()

		name := c.String("name")
		if name == "" {
			return fmt.Errorf("name is required")
		}

		l, err := lws.ReadLWS()
		if err != nil {
			return err
		}

		k := lws.GetKeyByName(l, name)
		if k == nil {
			return fmt.Errorf("specified key %q is not in the letter writing set %q", name, l.Name)
		}

		ak, err := access.AccessKey(l.Name, k)
		if err != nil {
			return err
		}
		defer ak.Close()

		testSha256 := sha256.Sum256([]byte("test"))
		if _, err := ak.Sign(testSha256[:]); err != nil {
			return err
		}
		s.Info("Confirmed that the key can sign a message without error.")

		return nil
	},
}
