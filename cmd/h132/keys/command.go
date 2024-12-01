package keys

import (
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    "keys",
	Aliases: []string{"k", "key"},
	Usage:   "Manage keys",
	Subcommands: []*cli.Command{
		addCommand,
		listCommand,
		rmCommand,
		testCommand,
	},
}
