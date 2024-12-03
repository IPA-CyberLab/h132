package lws

import (
	"github.com/urfave/cli/v2"

	"github.com/IPA-CyberLab/h132/lws"
)

var statusCommand = &cli.Command{
	Name:    "status",
	Aliases: []string{"s", "show", "dump"},
	Usage:   "Show the status of the letter writing set",
	Action: func(c *cli.Context) error {
		l, err := lws.ReadLWS()
		if err != nil {
			return err
		}

		return lws.DumpStatus(l)
	},
}
