package keys

import (
	"fmt"

	"github.com/IPA-CyberLab/h132/cmd/h132/promptrm"
	"github.com/IPA-CyberLab/h132/lws"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var rmCommand = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "Remove a key",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "name of the key to be removed",
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

		if lws.GetKeyByName(l, name) == nil {
			return fmt.Errorf("specified key %q is not in the letter writing set %q", name, l.Name)
		}

		if err := promptrm.Prompt(name, "key"); err != nil {
			return err
		}

		for i, key := range l.Keys {
			if key.Name == name {
				l.Keys = append(l.Keys[:i], l.Keys[i+1:]...)
				break
			}
		}

		if err := lws.UpdateLWS(l, lws.UpdateRemoveKey); err != nil {
			return fmt.Errorf("failed to update lws: %w", err)
		}

		s.Infof("Successfully removed the key %q from the letter writing set %q.", name, l.Name)
		return nil
	},
}
