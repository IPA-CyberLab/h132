package keys

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/IPA-CyberLab/h132/lws"
	"github.com/IPA-CyberLab/h132/pb"
)

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List keys",
	Action: func(c *cli.Context) error {
		s := zap.S()

		l, err := lws.ReadLWS()
		if err != nil {
			return err
		}

		s.Infof("Found %d keys in the letter writing set %q", len(l.Keys), l.Name)
		for i, key := range l.Keys {
			fmt.Printf("[%d] %s\n", i, pb.KeyImplSummary(key))
		}
		return nil
	},
}
